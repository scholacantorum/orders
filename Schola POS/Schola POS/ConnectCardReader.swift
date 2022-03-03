//
//  ConnectCardReader.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-28.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

class ConnectCardReader: UIViewController, CardReaderConnectDelegate {

    var reader: Reader!
    var completion: ChooseCardReaderDelegate!
    var update: ReaderSoftwareUpdate!
    var statusLabel: UILabel!
    var descLabel: UILabel!
    var applyButton: UIButton!
    var skipButton: UIButton!

    init(reader: Reader, completion: ChooseCardReaderDelegate) {
        super.init(nibName: nil, bundle: nil)
        self.reader = reader
        self.completion = completion
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

        cardReaderHandler.connect(reader, self)
    }

    func onUpdatingReader(_ update: ReaderSoftwareUpdate, progress: Float, canCancel: Bool) {
        statusLabel.text = "Applying required update to \(update.deviceSoftwareVersion), \(Int(progress*100.0))% complete"
    }
    
    func onReaderConnected() {
        self.dismiss(animated: false, completion: {
            self.completion.cardReaderReady()
        })
    }
    
    func onReaderConnectError(_ error: Error) {
        print("Error connecting to card reader: \(error)")
        let alert = UIAlertController(title: "Card Reader Error", message: error.localizedDescription, preferredStyle: .alert)
        alert.addAction(UIAlertAction(title: "OK", style: .default) { action in
            self.dismiss(animated: false, completion: nil)
        })
        self.present(alert, animated: true)
    }
}
