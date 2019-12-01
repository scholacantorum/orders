//
//  ChooseCardReader.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-28.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

// For testing, to skip login page:
// import Stripe

class ReaderCell: UITableViewCell {
    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: .subtitle, reuseIdentifier: reuseIdentifier)
    }
    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
}

protocol ChooseCardReaderDelegate {
    func cardReaderReady()
}

class ChooseCardReader: UIViewController, DiscoveryDelegate, UITableViewDataSource, UITableViewDelegate, ChooseCardReaderDelegate {

    var tableView: UITableView!
    var spinner: UIActivityIndicatorView!
    var discoverCancel: Cancelable?
    var readers = [Reader]()

    override func viewDidLoad() {
        super.viewDidLoad()
        navigationItem.hidesBackButton = true
        navigationItem.title = "Connect Card Reader"
        view.backgroundColor = UIColor.white
        backend.noop() // make sure it is initialized before using the Terminal

        // For testing, skip login page:
        // store.username = "sroth"
        // store.auth = "9999-9999-9999"
        // store.baseURL = "https://orders-test.scholacantorum.org/api"
        // store.allow.card = true
        // store.allow.cash = true
        // store.allow.willcall = true
        // STPPaymentConfiguration.shared().publishableKey = "pk_test_QPwvhWbGaakWn7DGcco8J5Nd"
        // store.event = Event(id: "2019-07-29", name: "Summer Sing", start: "2019-07-29T19:30:00", freeEntries: ["Student"])
        // store.products = [
        //     Product(id: "ticket-2019-07-29", name: "July 29", message: nil, price: 1700, ticketCount: 1),
        //     Product(id: "summer-sings-2019", name: "Flex Pass", message: nil, price: 8500, ticketCount: 6),
        //     Product(id: "summer-sings-2019-student", name: "Student", message: nil, price: 0, ticketCount: 1),
        // ]

        var constraints = [NSLayoutConstraint]()

        let spinnerHBox1 = UIView()
        spinnerHBox1.backgroundColor = UIColor.lightGray
        view.addSubview(spinnerHBox1)
        constraints.append(contentsOf: [
            spinnerHBox1.leftAnchor.constraint(equalTo: view.leftAnchor),
            spinnerHBox1.rightAnchor.constraint(equalTo: view.rightAnchor),
            spinnerHBox1.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
        ])

        let spinnerHBox2 = UIView()
        spinnerHBox1.addSubview(spinnerHBox2)
        constraints.append(contentsOf: [
            spinnerHBox2.centerXAnchor.constraint(equalTo: spinnerHBox1.centerXAnchor),
            spinnerHBox2.topAnchor.constraint(equalTo: spinnerHBox1.topAnchor, constant: 9.0),
            {
                let constraint = spinnerHBox2.heightAnchor.constraint(equalToConstant: 1.0)
                constraint.priority = .defaultLow
                return constraint
            }(),
            spinnerHBox2.bottomAnchor.constraint(equalTo: spinnerHBox1.bottomAnchor, constant: -9.0),
        ])

        let spinnerLabel = UILabel()
        spinnerLabel.text = "Looking for readers... "
        spinnerLabel.setContentHuggingPriority(.required, for: .horizontal)
        spinnerLabel.setContentCompressionResistancePriority(.required, for: .horizontal)
        spinnerHBox2.addSubview(spinnerLabel)
        constraints.append(contentsOf: [
            spinnerLabel.leftAnchor.constraint(equalTo: spinnerHBox2.leftAnchor),
            spinnerLabel.topAnchor.constraint(greaterThanOrEqualTo: spinnerHBox2.topAnchor),
            spinnerLabel.centerYAnchor.constraint(equalTo: spinnerHBox2.centerYAnchor),
            spinnerLabel.bottomAnchor.constraint(lessThanOrEqualTo: spinnerHBox2.bottomAnchor),
        ])

        spinner = UIActivityIndicatorView()
        spinner.style = UIActivityIndicatorView.Style.white
        spinner.color = UIColor.black
        spinner.setContentHuggingPriority(.required, for: .horizontal)
        spinner.setContentCompressionResistancePriority(.required, for: .horizontal)
        spinnerHBox2.addSubview(spinner)
        constraints.append(contentsOf: [
            spinner.leftAnchor.constraint(equalTo: spinnerLabel.rightAnchor),
            spinner.rightAnchor.constraint(equalTo: spinnerHBox2.rightAnchor),
            spinner.topAnchor.constraint(greaterThanOrEqualTo: spinnerHBox2.topAnchor),
            spinner.centerYAnchor.constraint(equalTo: spinnerHBox2.centerYAnchor),
            spinner.bottomAnchor.constraint(lessThanOrEqualTo: spinnerHBox2.bottomAnchor),
        ])

        tableView = UITableView()
        tableView.register(ReaderCell.self, forCellReuseIdentifier: "readerCell")
        tableView.delegate = self
        tableView.dataSource = self
        view.addSubview(tableView)
        constraints.append(contentsOf: [
            tableView.leftAnchor.constraint(equalTo: view.leftAnchor),
            tableView.rightAnchor.constraint(equalTo: view.rightAnchor),
            tableView.topAnchor.constraint(equalTo: spinnerHBox1.bottomAnchor),
        ])

        let skipButton = UIButton()
        skipButton.setTitle("No Card Reader", for: .normal)
        skipButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        skipButton.setTitleColor(UIColor.white, for: .normal)
        skipButton.backgroundColor = UIColor.darkGray
        skipButton.layer.cornerRadius = 5.0
        skipButton.addTarget(self, action: #selector(skipButton(_:)), for: .touchUpInside)
        skipButton.setContentHuggingPriority(.required, for: .vertical)
        skipButton.setContentCompressionResistancePriority(.required, for: .vertical)
        view.addSubview(skipButton)
        constraints.append(contentsOf: [
            skipButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            skipButton.widthAnchor.constraint(equalToConstant: 180.0),
            skipButton.topAnchor.constraint(equalTo: tableView.bottomAnchor, constant: 9.0),
            skipButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])

        NSLayoutConstraint.useAndActivateConstraints(constraints)
    }

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return readers.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "readerCell", for: indexPath)
        let reader = readers[indexPath.row]
        cell.textLabel!.text = "\(reader.simulated ? "Simulated " : "")\(Terminal.stringFromDeviceType(reader.deviceType)) \(reader.serialNumber)"
        var batteryLevelString: String
        if let bl = reader.batteryLevel {
            batteryLevelString = "\(Int(Double(truncating: bl)*100.0))%"
        } else {
            batteryLevelString = "unknown"
        }
        cell.detailTextLabel!.text = "Battery: \(batteryLevelString)     Version: \(reader.deviceSoftwareVersion ?? "unknown")"
        return cell
    }

    func connect() {
        if Terminal.shared.connectionStatus != .notConnected {
            navigationController!.pushViewController(Main(), animated: true)
            return
        }
        spinner.startAnimating()
        discoverCancel = Terminal.shared.discoverReaders(DiscoveryConfiguration(deviceType: .chipper2X, discoveryMethod: .bluetoothScan, simulated: false), delegate: self) { error in
            if let error = error {
                print("Error discovering card readers: \(error)")
                let alert = UIAlertController(title: "Card Reader Error", message: error.localizedDescription, preferredStyle: .alert)
                alert.addAction(UIAlertAction(title: "OK", style: .default))
                self.present(alert, animated: true, completion: nil)
                self.spinner.stopAnimating()
                return
            }
        }
    }

    func terminal(_ terminal: Terminal, didUpdateDiscoveredReaders readers: [Reader]) {
        self.readers = readers
        tableView.reloadData()
    }

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        let reader = readers[indexPath.row]
        present(ConnectCardReader(reader: reader, completion: self), animated: true, completion: nil)
    }

    @objc func skipButton(_ sender: UIButton) {
        spinner.stopAnimating()
        if let cancel = discoverCancel {
            if !cancel.completed {
                cancel.cancel { error in
                    if let error = error {
                        print("Error canceling card reader discovery: \(error)")
                    }
                }
            }
        }
        navigationController!.pushViewController(Main(), animated: true)
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)
        connect()
    }

    func cardReaderReady() {
        navigationController!.pushViewController(Main(), animated: true)
    }

}
