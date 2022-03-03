//
//  CardReaderHandler.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//
// This class handles all of the mechanics of discovering, connecting to, and managing
// the Bluetooth card reader.

import StripeTerminal

// The cardReaderHandler is the singleton instance of CardReaderHandler.
let cardReaderHandler = CardReaderHandler()

// CardReaderDiscoveryDelegate is the delegate passed into the discover() method.  It
// must remain valid until connect() is called, cancelDiscover() is called, or
// onReaderDiscoveryError() is called.
protocol CardReaderDiscoveryDelegate {
    // onAvailableReaders is called when the list of discovered readers changes.
    func onAvailableReaders(_ readers: [Reader])
    // onReaderDiscoveryError is called when an error occurs during discovery.
    func onReaderDiscoveryError(_ error: Error)
}

// CardReaderConnectDelegate is the delegate passed into the connect() method.  It
// must remain valid until onReaderConnected() or onReaderConnectError() is called.
protocol CardReaderConnectDelegate {
    // onUpdatingReader is called when a required update is being installed.  It
    // may be called zero or more times, with progress values increasing
    // monotonically from zero to one.  It is never called after onReaderConnected
    // or onReaderConnectError.
    func onUpdatingReader(_ update: ReaderSoftwareUpdate, progress: Float, canCancel: Bool)
    // onReaderConnected is called when the connection is successfully established.
    func onReaderConnected()
    // onReaderConnectError is called when the connection to the reader fails.
    func onReaderConnectError(_ error: Error)
}

struct CardReaderStatus {
    var connectionStatus: ConnectionStatus
    var batteryLevel: Float
    var batteryStatus: BatteryStatus
    var batteryCharging: Bool
    var updateAvailable: ReaderSoftwareUpdate?
}

// CardReaderStatusDelegate is a protocol honored by a component that needs to be
// informed of changes in card reader status or battery level.  There can
// be only one such component at a time.  The status delegate can be set with the
// cardReaderHandler.setStatusDelegate method.
protocol CardReaderStatusDelegate {
    // onCardReaderStatusChange is called when the card reader status changes.
    func onCardReaderStatusChange(_ status: CardReaderStatus)
}

// CardReaderDisplayDelegate handles presenting messages to the consumer for the
// reader.
protocol CardReaderDisplayDelegate {
    // onRequestReaderInput is called to request the customer to present a payment
    // method to the reader.
    func onRequestReaderInput(_ options: String)
    // onDisplayMessage is called to request the customer take an action with the
    // reader.
    func onDisplayMessage(_ message: String)
}

// CardReaderUpdateDelegate handles optional software updates to the reader.
protocol CardReaderUpdateDelegate {
    // onUpdatingReader is called when an option update is being installed.  It
    // will be called multiple times, with progress values increasing monotonically
    // from zero to one.  It is never called after onUpdateComplete.
    func onUpdatingReader(_ update: ReaderSoftwareUpdate, progress: Float, canCancel: Bool)
    // onUpdateCompleted is called when the update has finished, successfully or not.
    func onUpdateComplete(_ error: Error?)
}

class CardReaderHandler: NSObject, TerminalDelegate, DiscoveryDelegate, BluetoothReaderDelegate {

    var discoveryDelegate: CardReaderDiscoveryDelegate?
    var connectionDelegate: CardReaderConnectDelegate?
    var statusDelegate: CardReaderStatusDelegate?
    var displayDelegate: CardReaderDisplayDelegate?
    var updateDelegate: CardReaderUpdateDelegate?
    var status: CardReaderStatus
    var discoverCancel: Cancelable?
    var updateInProgress: ReaderSoftwareUpdate?
    var updateCancel: Cancelable?

    override init() {
        status = CardReaderStatus(
            connectionStatus: Terminal.shared.connectionStatus,
            batteryLevel: 0,
            batteryStatus: BatteryStatus.unknown,
            batteryCharging: false,
            updateAvailable: nil
        )
        super.init()
        Terminal.shared.delegate = self
    }
    
    // MARK:  Reader Discovery
    
    // discover starts the process of discovering a card reader.  The first callback
    // function is called whenever the list of available readers changes.  The
    // second callback function is called (and discovery is aborted) if any error
    // occurs.  Discovery continues until connection to a discovered reader is requested,
    // an error occurs, or cancelDiscover is called.
    func discover(_ delegate: CardReaderDiscoveryDelegate) {
        discoveryDelegate = delegate
        let config = DiscoveryConfiguration(discoveryMethod: .bluetoothScan, simulated: false)
        discoverCancel = Terminal.shared.discoverReaders(config, delegate: self) { error in
            if let error = error {
                print("Error discovering card readers: \(error)")
                delegate.onReaderDiscoveryError(error)
                self.discoverCancel = nil
                return
            }
        }
    }
    
    // terminal:didUpdateDiscoveredReaders is called when the list of discovered
    // readers changes.  It passes them on to the callback provided to discover().
    func terminal(_ terminal: Terminal, didUpdateDiscoveredReaders readers: [Reader]) {
        discoveryDelegate?.onAvailableReaders(readers)
    }

    // cancelDiscover terminates any ongoing reader discovery operation.
    func cancelDiscover() {
        discoverCancel?.cancel() { error in
            if let error = error {
                print("Error canceling reader discovery: \(error)")
                self.discoveryDelegate?.onReaderDiscoveryError(error)
            }
            self.discoverCancel = nil
            self.discoveryDelegate = nil
        }
    }
    
    // MARK:  Reader Connection
    
    func connect(_ reader: Reader, _ delegate: CardReaderConnectDelegate) {
        discoveryDelegate = nil
        connectionDelegate = delegate
        var config: BluetoothConnectionConfiguration
        if let loc = reader.locationId {
            config = BluetoothConnectionConfiguration(locationId: loc)
        } else {
            config = BluetoothConnectionConfiguration(locationId: store.testmode ? "tml_EPZVrQJ9viIwM1" : "tml_EPZV9AjTYoLTIs")
        }
        Terminal.shared.connectBluetoothReader(reader, delegate: self, connectionConfig: config) { reader, error in
            if let error = error {
                print("Error connecting to card reader: \(error)")
                delegate.onReaderConnectError(error)
            } else {
                delegate.onReaderConnected()
            }
            self.connectionDelegate = nil
        }
    }
    
    
    func disconnect() {
        if status.connectionStatus == .connected {
            Terminal.shared.disconnectReader() { error in
                if let error = error {
                    print("Error disconnecting card reader: \(error)")
                }
            }
        }
    }

    // MARK:  Reader Status Updates
    
    // setStatusDelegate sets the card reader status delegate (or clears it, if the
    // argument is nil).  The new delegate will be called immediately, and whenever
    // the card reader status changes or the card reader raises a low battery alert.
    func setStatusDelegate(_ delegate: CardReaderStatusDelegate?) {
        self.statusDelegate = delegate
        delegate?.onCardReaderStatusChange(status)
    }

    func terminal(_ terminal: Terminal, didChangeConnectionStatus cstat: ConnectionStatus) {
        status.connectionStatus = cstat
        statusDelegate?.onCardReaderStatusChange(status)
    }

    func terminal(_ terminal: Terminal, didReportUnexpectedReaderDisconnect reader: Reader) {
        print("Unexpected reader disconnect")
    }

    func reader(_ reader: Reader, didReportAvailableUpdate update: ReaderSoftwareUpdate) {
        status.updateAvailable = update
        statusDelegate?.onCardReaderStatusChange(status)
    }
    
    func readerDidReportLowBatteryWarning(_ reader: Reader) {
        print("Card reader low battery")
        if status.batteryStatus != BatteryStatus.critical {
            status.batteryStatus = BatteryStatus.low
        }
        statusDelegate?.onCardReaderStatusChange(status)
    }

    func reader(_ reader: Reader, didReportBatteryLevel batteryLevel: Float, status bstat: BatteryStatus, isCharging: Bool) {
        status.batteryLevel = batteryLevel
        status.batteryStatus = bstat
        status.batteryCharging = isCharging
        statusDelegate?.onCardReaderStatusChange(status)
    }
   
    // MARK:  Reader Display Updates

    // setDisplayDelegate sets the card reader display delegate (or clears it, if the
    // argument is nil).
    func setDisplayDelegate(_ delegate: CardReaderDisplayDelegate?) {
        displayDelegate = delegate
    }
    
    func reader(_ reader: Reader, didRequestReaderInput inputOptions: ReaderInputOptions = []) {
        if let delegate = displayDelegate {
            delegate.onRequestReaderInput(Terminal.stringFromReaderInputOptions(inputOptions))
        } else {
            print("Unexpected request for reader input")
        }
    }
    
    func reader(_ reader: Reader, didRequestReaderDisplayMessage displayMessage: ReaderDisplayMessage) {
        if let delegate = displayDelegate {
            delegate.onDisplayMessage(Terminal.stringFromReaderDisplayMessage(displayMessage))
        } else {
            print("Unexpected request for reader display")
        }
    }
    
    // MARK:  Reader Software Updates
    
    // installUpdate starts the installation of the available update.
    func installUpdate(_ delegate: CardReaderUpdateDelegate) {
        updateDelegate = delegate
        Terminal.shared.installAvailableUpdate()
    }
    
    func reader(_ reader: Reader, didStartInstallingUpdate update: ReaderSoftwareUpdate, cancelable: Cancelable?) {
        updateInProgress = update
        updateCancel = cancelable
        connectionDelegate?.onUpdatingReader(update, progress: 0, canCancel: cancelable != nil)
        updateDelegate?.onUpdatingReader(update, progress: 0, canCancel: cancelable != nil)
    }
    
    func reader(_ reader: Reader, didReportReaderSoftwareUpdateProgress progress: Float) {
        connectionDelegate?.onUpdatingReader(updateInProgress!, progress: progress, canCancel: updateCancel != nil)
        updateDelegate?.onUpdatingReader(updateInProgress!, progress: progress, canCancel: updateCancel != nil)
    }
    
    func reader(_ reader: Reader, didFinishInstallingUpdate update: ReaderSoftwareUpdate?, error: Error?) {
        if let error = error {
            connectionDelegate?.onReaderConnectError(error)
        } else {
            connectionDelegate?.onUpdatingReader(updateInProgress!, progress: 1.0, canCancel: false)
        }
        updateDelegate?.onUpdateComplete(error)
        updateInProgress = nil
        updateCancel = nil
        updateDelegate = nil
    }

    func cancelUpdate(_ completion: @escaping ErrorCompletionBlock) {
        updateCancel?.cancel(completion)
        updateCancel = nil
        updateInProgress = nil
        updateDelegate = nil
    }

}
