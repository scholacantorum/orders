//
//  CardReaderBanner.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

let warningColor = UIColor(red: 1.0, green: 0.8, blue: 0.4, alpha: 1.0)
let bannerHeight: CGFloat = 45.0

class CardReaderBanner: UIViewController, CardReaderStatusDelegate {

    var reconnect: () -> Void
    var status = ConnectionStatus.notConnected
    var lowBattery = false

    lazy var stateLabel = UILabel()
    var heightConstraint: NSLayoutConstraint!

    init(reconnect: @escaping () -> Void) {
        self.reconnect = reconnect
        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        stateLabel.textAlignment = .center
        view.addSubview(stateLabel)
        heightConstraint = stateLabel.heightAnchor.constraint(equalToConstant: bannerHeight)
        NSLayoutConstraint.useAndActivateConstraints([
            stateLabel.leftAnchor.constraint(equalTo: view.leftAnchor),
            stateLabel.rightAnchor.constraint(equalTo: view.rightAnchor),
            stateLabel.topAnchor.constraint(equalTo: view.topAnchor),
            stateLabel.bottomAnchor.constraint(equalTo: view.bottomAnchor),
            heightConstraint,
        ])
        cardReaderHandler.setStatusDelegate(self)
        view.addGestureRecognizer(UITapGestureRecognizer(target: self, action: #selector(handleTap(_:))))
    }

    func onCardReaderStatusChange(_ status: CardReaderStatus) {
        self.status = status.connectionStatus
        self.lowBattery = status.batteryStatus == BatteryStatus.low || status.batteryStatus == BatteryStatus.critical
        setStateLabel()
    }

    func setStateLabel() {
        switch status {
        case .connecting:
            heightConstraint.constant = bannerHeight
            stateLabel.text = "Connecting to card reader"
            stateLabel.backgroundColor = warningColor
        case .notConnected:
            heightConstraint.constant = bannerHeight
            stateLabel.text = "Card reader not connected"
            stateLabel.backgroundColor = UIColor.red
        default:
            if lowBattery {
                heightConstraint.constant = bannerHeight
                stateLabel.text = "Card reader has low battery"
                stateLabel.backgroundColor = warningColor
            } else {
                heightConstraint.constant = 0.0
                stateLabel.text = nil
            }
        }
    }

    @objc func handleTap(_ sender: UITapGestureRecognizer) {
        if status != .connected {
            reconnect()
        }
    }
    
}
