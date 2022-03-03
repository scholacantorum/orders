//
//  SalesPaymentCardPresent.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

class SalesPaymentCardPresent: UIViewController, CardReaderDisplayDelegate {
    

    var order: Order!
    var intent: PaymentIntent!
    var collectCancel: Cancelable?

    lazy var messageLabel: UILabel = {
        let view = UILabel()
        view.numberOfLines = 0
        view.lineBreakMode = .byWordWrapping
        view.textAlignment = .center
        view.font = UIFont.boldSystemFont(ofSize: 30.0)
        view.textColor = scholaBlue
        view.text = "This is a test of the emergency broadcasting system."
        return view
    }()
    lazy var manualEntryButton: UIButton = {
        let view = UIButton()
        view.setTitle("Manual Entry", for: .normal)
        view.setContentHuggingPriority(.required, for: .vertical)
        view.setContentCompressionResistancePriority(.required, for: .vertical)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(manualEntryButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var cancelButton: UIButton = {
        let view = UIButton()
        view.setTitle("Cancel", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = UIColor.darkGray
        view.addTarget(self, action: #selector(cancelButton(_:)), for: .touchUpInside)
        return view
    }()

    init(order: Order) {
        super.init(nibName: nil, bundle: nil)
        self.order = order
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white

        let summaryController = SalesSummary(order: order)
        let summaryView = UIView()
        summaryController.view = summaryView
        summaryController.viewDidLoad()
        addChild(summaryController)
        view.addSubview(summaryView)
        view.addSubview(messageLabel)
        view.addSubview(manualEntryButton)
        view.addSubview(cancelButton)
        NSLayoutConstraint.useAndActivateConstraints([
            summaryView.leftAnchor.constraint(equalTo: view.leftAnchor),
            summaryView.rightAnchor.constraint(equalTo: view.rightAnchor),
            summaryView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 18.0),
            {
                let constraint = summaryView.heightAnchor.constraint(equalToConstant: 1.0)
                constraint.priority = .defaultLow
                return constraint
            }(),
            messageLabel.leftAnchor.constraint(equalTo: view.leftAnchor, constant: 9.0),
            messageLabel.rightAnchor.constraint(equalTo: view.rightAnchor, constant: -9.0),
            messageLabel.topAnchor.constraint(equalTo: summaryView.bottomAnchor),
            messageLabel.bottomAnchor.constraint(equalTo: manualEntryButton.topAnchor, constant: -9.0),
            manualEntryButton.widthAnchor.constraint(equalToConstant: 150.0),
            manualEntryButton.rightAnchor.constraint(equalTo: view.centerXAnchor, constant: -9.0),
            manualEntryButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
            cancelButton.leftAnchor.constraint(equalTo: view.centerXAnchor, constant: 9.0),
            cancelButton.widthAnchor.constraint(equalToConstant: 150.0),
            cancelButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])

        createOrder()
    }

    func createOrder() {
        messageLabel.text = "Creating order..."
        messageLabel.textColor = scholaBlue
        backend.placeOrder(order: order) { order, error in
            DispatchQueue.main.async {
                if let error = error {
                    self.messageLabel.text = error
                    self.messageLabel.textColor = UIColor.red
                    return
                }
                self.order = order
                self.retrieveIntent()
            }
        }
    }

    func retrieveIntent() {
        Terminal.shared.retrievePaymentIntent(clientSecret: order.payments[0].method!) { intent, error in
            DispatchQueue.main.async {
                if let error = error {
                    self.messageLabel.text = error.localizedDescription
                    self.messageLabel.textColor = UIColor.red
                    backend.cancelOrder(order: self.order)
                    return
                }
                self.intent = intent
                self.collectPayment()
            }
        }
    }

    func collectPayment() {
        cardReaderHandler.setDisplayDelegate(self)
        collectCancel = Terminal.shared.collectPaymentMethod(intent) { intent, error in
            cardReaderHandler.setDisplayDelegate(nil)
            DispatchQueue.main.async {
                if let error = error as NSError? {
                    if error.code == 2020 {
                        return
                    }
                    self.messageLabel.text = error.localizedDescription
                    self.messageLabel.textColor = UIColor.red
                    return
                }
                self.intent = intent
                self.processPayment()
            }
        }
    }

    func processPayment() {
        messageLabel.text = "Processing payment..."
        messageLabel.textColor = scholaBlue
        Terminal.shared.processPayment(intent) { intent, error in
            DispatchQueue.main.async {
                if let error = error as ProcessPaymentError? {
                    if error.paymentIntent == nil || error.paymentIntent?.status == .requiresConfirmation {
                        self.messageLabel.text = error.requestError?.localizedDescription
                        self.messageLabel.textColor = UIColor.red
                        return
                    }
                    let alert = UIAlertController(title: "Payment Error", message: error.localizedDescription, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true) {
                        self.collectPayment()
                    }
                    return
                }
                self.manualEntryButton.isEnabled = false
                self.manualEntryButton.backgroundColor = scholaBlueDisabled
                self.cancelButton.isEnabled = false
                self.cancelButton.backgroundColor = UIColor.lightGray
                self.intent = intent
                for line in self.order.lines {
                    store.admitted += line.used ?? 0
                    if line.product != "donation" {
                        store.sold += line.quantity
                    }
                }
                self.capturePayment()
            }
        }
    }

    func capturePayment() {
        messageLabel.text = "Recording successful payment..."
        messageLabel.textColor = scholaBlue
        backend.captureOrderPayment(order: order) { order, error in
            DispatchQueue.main.async {
                if error != nil {
                    // We are not displaying an error here because the credit card did get
                    // authorized, and we'll eventually notice the authorized but uncaptured
                    // payment and capture it.  It's a problem, but not one the person doing
                    // at-the-door sales should be distracted with.  However, we'll roll
                    // back to the main menu rather than offering to send a receipt.
                    self.navigationController!.popToRootViewController(animated: true)
                    return
                }
                if self.order.email ?? "" != "" {
                    // Already have an email address, so already sent a receipt.
                    self.navigationController?.popToRootViewController(animated: true)
                    return
                }
                self.navigationController!.pushViewController(SalesReceipt(order: order!), animated: true)
            }
        }
    }

    func onRequestReaderInput(_ options: String) {
        self.messageLabel.text = options
        self.messageLabel.textColor = warningColor
    }
    
    func onDisplayMessage(_ message: String) {
        self.messageLabel.text = message
        self.messageLabel.textColor = warningColor
    }

    @objc func manualEntryButton(_ sender: UIButton) {
        cancelInner({
            self.order.id = nil
            self.order.payments[0].type = "card"
            let nc = self.navigationController!
            nc.popViewController(animated: false)
            nc.pushViewController(SalesPaymentCard(order: self.order), animated: true)
        })
    }

    @objc func cancelButton(_ sender: UIButton) {
        cancelInner({ self.navigationController!.popViewController(animated: true) })
    }

    func cancelInner(_ handler: @escaping () -> Void) {
        if let intent = intent {
            if intent.status == .requiresCapture {
                return
            }
            backend.cancelOrder(order: order)
        }
        if let cancel = collectCancel {
            if !cancel.completed {
                cancel.cancel { error in
                    DispatchQueue.main.async {
                        if let error = error {
                            print("Error canceling collection of payment: \(error.localizedDescription)")
                        }
                        handler()
                    }
                }
                return
            }
        }
        handler()
    }
}
