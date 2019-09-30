//
//  MainMenuController.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-24.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

class MainMenu: UIViewController {

    var reconnect: () -> Void

    lazy var sellTicketsButton: UIButton = {
        let view = UIButton()
        view.setTitle("Sell Tickets", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 30.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(sellTicketsButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var willCallButton: UIButton = {
        let view = UIButton()
        view.setTitle("Will Call", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 30.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(willCallButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var scanTicketButton: UIButton = {
        let view = UIButton()
        view.setTitle("Scan Ticket", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 30.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(scanTicketButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var sellMerchandiseButton: UIButton = {
        let view = UIButton()
        view.setTitle("Sell Merchandise", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 30.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(sellMerchandiseButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var logoutButton: UIButton = {
        let view = UIButton()
        view.setTitle("Logout", for: .normal)
        view.titleLabel!.font = UIFont.systemFont(ofSize: 18.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = UIColor.darkGray
        view.addTarget(self, action: #selector(logoutButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var statsLabel: UILabel = {
        let view = UILabel()
        view.textColor = UIColor.darkGray
        view.font = UIFont.systemFont(ofSize: 15.0)
        view.numberOfLines = 0
        view.textAlignment = .center
        return view
    }()

    init(reconnect: @escaping () -> Void) {
        self.reconnect = reconnect
        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white

        var constraints: [NSLayoutConstraint] = []
        var topAnchor = view.safeAreaLayoutGuide.topAnchor

        if store.allow.card {
            let cardReaderView = UIView()
            let cardReaderBanner = CardReaderBanner(reconnect: reconnect)
            cardReaderBanner.view = cardReaderView
            cardReaderBanner.viewDidLoad()
            addChild(cardReaderBanner)
            view.addSubview(cardReaderView)
            constraints.append(contentsOf: [
                cardReaderView.leftAnchor.constraint(equalTo: view.leftAnchor),
                cardReaderView.rightAnchor.constraint(equalTo: view.rightAnchor),
                cardReaderView.topAnchor.constraint(equalTo: topAnchor),
            ])
            topAnchor = cardReaderView.bottomAnchor
        }
        if store.allow.card || store.allow.cash {
            view.addSubview(sellTicketsButton)
            constraints.append(contentsOf: [
                sellTicketsButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
                sellTicketsButton.topAnchor.constraint(equalTo: topAnchor, constant: 45.0),
                sellTicketsButton.widthAnchor.constraint(equalTo: scanTicketButton.widthAnchor),
            ])
            topAnchor = sellTicketsButton.bottomAnchor
        }
        if store.allow.willcall {
            view.addSubview(willCallButton)
            constraints.append(contentsOf: [
                willCallButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
                willCallButton.topAnchor.constraint(equalTo: topAnchor, constant: 45.0),
                willCallButton.widthAnchor.constraint(equalTo: scanTicketButton.widthAnchor),
            ])
            topAnchor = willCallButton.bottomAnchor
        }
        view.addSubview(scanTicketButton)
//        view.addSubview(sellMerchandiseButton)
        view.addSubview(logoutButton)
        view.addSubview(statsLabel)
        constraints.append(contentsOf: [
            scanTicketButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            scanTicketButton.topAnchor.constraint(equalTo: topAnchor, constant: 45.0),
            scanTicketButton.widthAnchor.constraint(equalTo: scanTicketButton.titleLabel!.widthAnchor, constant: 36.0),
//            sellMerchandiseButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
//            sellMerchandiseButton.topAnchor.constraint(equalTo: scanTicketButton.bottomAnchor, constant: 45.0),
//            sellMerchandiseButton.widthAnchor.constraint(equalTo: sellMerchandiseButton.titleLabel!.widthAnchor, constant: 36.0),
            logoutButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            logoutButton.topAnchor.constraint(equalTo: scanTicketButton.bottomAnchor, constant: 60.0),
//            logoutButton.topAnchor.constraint(equalTo: sellMerchandiseButton.bottomAnchor, constant: 60.0),
            logoutButton.widthAnchor.constraint(equalTo: logoutButton.titleLabel!.widthAnchor, constant: 18.0),
            statsLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            statsLabel.topAnchor.constraint(equalTo: logoutButton.bottomAnchor, constant: 9.0),
        ])
        NSLayoutConstraint.useAndActivateConstraints(constraints)
        statsLabel.text = "Statistics here"
    }

    override func viewWillAppear(_ animated: Bool) {
        statsLabel.text = "Admitted \(store.admitted)\nSold \(store.sold)\nCash $\(store.cash / 100)\nCheck $\(store.check / 100)"
        if store.allow.card {
            if let reader = Terminal.shared.connectedReader {
                if let battery = reader.batteryLevel {
                    statsLabel.text = statsLabel.text! + "\nReader battery \(Int(Double(truncating: battery)*100.0))%"
                }
            }
        }
    }

    @objc func sellTicketsButton(_ sender: UIButton) {
        navigationController!.pushViewController(SellTickets(), animated: true)
    }

    @objc func willCallButton(_ sender: UIButton) {
        navigationController!.pushViewController(WillCall(), animated: true)
    }

    @objc func scanTicketButton(_ sender: UIButton) {
        navigationController!.pushViewController(ScanTicket(), animated: true)
    }

    @objc func sellMerchandiseButton(_ sender: UIButton) {
        navigationController!.pushViewController(SellMerchandise(), animated: true)
    }

    @objc func logoutButton(_ sender: UIButton) {
        store.logout()
    }

}
