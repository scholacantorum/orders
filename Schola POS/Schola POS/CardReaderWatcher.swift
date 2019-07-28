//
//  CardReaderWatcher.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import StripeTerminal

let cardReaderWatcher = CardReaderWatcher()

protocol CardReaderWatcherDelegate {
    func cardReader(status: ConnectionStatus, lowBattery: Bool)
}

class CardReaderWatcher: NSObject, TerminalDelegate {

    var delegate: CardReaderWatcherDelegate?
    var lowBattery = false

    func connect() {
        Terminal.shared.delegate = self
    }

    func setDelegate(_ delegate: CardReaderWatcherDelegate?) {
        self.delegate = delegate
        delegate?.cardReader(status: Terminal.shared.connectionStatus, lowBattery: lowBattery)
    }

    func terminal(_ terminal: Terminal, didChangeConnectionStatus status: ConnectionStatus) {
        delegate?.cardReader(status: Terminal.shared.connectionStatus, lowBattery: lowBattery)
    }

    func terminal(_ terminal: Terminal, didReportUnexpectedReaderDisconnect reader: Reader) {
        print("Unexpected reader disconnect")
    }

    func terminalDidReportLowBatteryWarning(_ terminal: Terminal) {
        print("Card reader low battery")
        lowBattery = true
        delegate?.cardReader(status: Terminal.shared.connectionStatus, lowBattery: lowBattery)
    }

    func disconnect() {
        if Terminal.shared.delegate != nil {
            if Terminal.shared.connectionStatus == .connected {
                Terminal.shared.disconnectReader() { error in
                    if let error = error {
                        print("Error disconnecting card reader: \(error)")
                    }
                }
            }
        }
    }

}
