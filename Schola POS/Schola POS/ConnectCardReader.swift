//
//  ConnectCardReader.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-28.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

class ConnectCardReader: UIViewController, ReaderSoftwareUpdateDelegate {

    var reader: Reader!
    var update: ReaderSoftwareUpdate!
    var statusLabel: UILabel!
    var descLabel: UILabel!
    var applyButton: UIButton!
    var skipButton: UIButton!

    init(reader: Reader) {
        super.init(nibName: nil, bundle: nil)
        self.reader = reader
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white

        statusLabel = UILabel()
        statusLabel.text = "Connecting..."
        statusLabel.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.addSubview(statusLabel)
        NSLayoutConstraint.useAndActivateConstraints([
            statusLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            statusLabel.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 18.0),
        ])

        Terminal.shared.connectReader(reader) { reader, error in
            if let error = error {
                print("Error connecting to card reader: \(error)")
                let alert = UIAlertController(title: "Card Reader Error", message: error.localizedDescription, preferredStyle: .alert)
                alert.addAction(UIAlertAction(title: "OK", style: .default) { action in
                    self.dismiss(animated: false, completion: nil)
                })
                self.present(alert, animated: true)
                return
            }
            self.connected()
        }
    }

    func connected() {
        cardReaderWatcher.connect()
        statusLabel.text = "Checking for update..."
        Terminal.shared.checkForUpdate() { update, error in
            if let error = error {
                print("Error connecting to card reader: \(error)")
                let alert = UIAlertController(title: "Card Reader Error", message: error.localizedDescription, preferredStyle: .alert)
                alert.addAction(UIAlertAction(title: "OK", style: .default) { action in
                    self.dismiss(animated: false, completion: nil)
                })
                self.present(alert, animated: true)
                return
            }
            guard let update = update else {
                self.dismiss(animated: false, completion: nil)
                return
            }
            self.update = update
            self.askApplyUpdate()
        }
    }

    func askApplyUpdate() {
        statusLabel.text = "Card Reader Update"

        descLabel = UILabel()
        descLabel.text = "A card reader update is available.\nThe card reader currently has \(Terminal.shared.connectedReader!.deviceSoftwareVersion ?? "unknown").\nThe new version is \(update.deviceSoftwareVersion).\nThe update should take \(ReaderSoftwareUpdate.string(from: update.estimatedUpdateTime)).\nDo you want to apply the update now?"
        descLabel.lineBreakMode = .byWordWrapping
        descLabel.numberOfLines = 0
        view.addSubview(descLabel)

        applyButton = UIButton()
        applyButton.setTitle("Apply Update", for: .normal)
        applyButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        applyButton.setTitleColor(UIColor.white, for: .normal)
        applyButton.layer.cornerRadius = 5.0
        applyButton.backgroundColor = scholaBlue
        applyButton.addTarget(self, action: #selector(applyButton(_:)), for: .touchUpInside)
        view.addSubview(applyButton)

        skipButton = UIButton()
        skipButton.setTitle("Skip Update", for: .normal)
        skipButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        skipButton.setTitleColor(UIColor.white, for: .normal)
        skipButton.layer.cornerRadius = 5.0
        skipButton.backgroundColor = UIColor.darkGray
        skipButton.addTarget(self, action: #selector(skipButton(_:)), for: .touchUpInside)
        view.addSubview(skipButton)

        NSLayoutConstraint.useAndActivateConstraints([
            descLabel.leftAnchor.constraint(equalTo: view.leftAnchor, constant: 9.0),
            descLabel.rightAnchor.constraint(equalTo: view.rightAnchor, constant: 9.0),
            descLabel.topAnchor.constraint(equalTo: statusLabel.bottomAnchor, constant: 18.0),
            applyButton.widthAnchor.constraint(equalToConstant: 150.0),
            applyButton.rightAnchor.constraint(equalTo: view.centerXAnchor, constant: -9.0),
            applyButton.topAnchor.constraint(equalTo: descLabel.bottomAnchor, constant: 18.0),
            skipButton.leftAnchor.constraint(equalTo: view.centerXAnchor, constant: 9.0),
            skipButton.widthAnchor.constraint(equalToConstant: 150.0),
            skipButton.topAnchor.constraint(equalTo: descLabel.bottomAnchor, constant: 18.0),
        ])

    }

    @objc func applyButton(_ sender: UIButton) {
        descLabel.text = "Updating to \(update.deviceSoftwareVersion)..."
        applyButton.removeFromSuperview()
        skipButton.removeFromSuperview()
        Terminal.shared.installUpdate(update, delegate: self) { error in
            if let error = error {
                print("Error updating card reader: \(error)")
                let alert = UIAlertController(title: "Card Reader Error", message: error.localizedDescription, preferredStyle: .alert)
                alert.addAction(UIAlertAction(title: "OK", style: .default) { action in
                    self.dismiss(animated: false, completion: nil)
                })
                self.present(alert, animated: true)
                return
            }
            self.dismiss(animated: false, completion: nil)
        }
    }

    func terminal(_ terminal: Terminal, didReportReaderSoftwareUpdateProgress progress: Float) {
        descLabel.text = "Updating to \(update.deviceSoftwareVersion)...\n\(Int(progress*100.0))% complete."
    }

    @objc func skipButton(_ sender: UIButton) {
        self.dismiss(animated: false, completion: nil)
    }

}
