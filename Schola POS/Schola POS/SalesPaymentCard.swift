//
//  SalesPaymentCard.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import Stripe

class SalesPaymentCard: UIViewController, STPPaymentCardTextFieldDelegate {

    var order: Order!
    var cardEntry: STPPaymentCardTextField!
    var payNowButton: UIButton!

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

        let paymentCardLabel = UILabel()
        paymentCardLabel.text = "Payment Card"
        view.addSubview(paymentCardLabel)

        cardEntry = STPPaymentCardTextField()
        cardEntry.delegate = self
        view.addSubview(cardEntry)

        payNowButton = UIButton()
        payNowButton.setTitle("Pay Now", for: .normal)
        payNowButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        payNowButton.setTitleColor(UIColor.white, for: .normal)
        payNowButton.layer.cornerRadius = 5.0
        payNowButton.backgroundColor = scholaBlueDisabled
        payNowButton.isEnabled = false
        payNowButton.addTarget(self, action: #selector(payNowButton(_:)), for: .touchUpInside)
        view.addSubview(payNowButton)

        let cancelButton = UIButton()
        cancelButton.setTitle("Cancel", for: .normal)
        cancelButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        cancelButton.setTitleColor(UIColor.white, for: .normal)
        cancelButton.layer.cornerRadius = 5.0
        cancelButton.backgroundColor = UIColor.darkGray
        cancelButton.addTarget(self, action: #selector(cancelButton(_:)), for: .touchUpInside)
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
            paymentCardLabel.leftAnchor.constraint(equalTo: view.leftAnchor, constant: 9.0),
            paymentCardLabel.rightAnchor.constraint(equalTo: view.rightAnchor, constant: -9.0),
            paymentCardLabel.topAnchor.constraint(equalTo: summaryView.bottomAnchor, constant: 36.0),
            cardEntry.leftAnchor.constraint(equalTo: view.leftAnchor, constant: 9.0),
            cardEntry.rightAnchor.constraint(equalTo: view.rightAnchor, constant: -9.0),
            cardEntry.topAnchor.constraint(equalTo: paymentCardLabel.bottomAnchor),
            cardEntry.heightAnchor.constraint(equalToConstant: 44.0),
            payNowButton.widthAnchor.constraint(equalToConstant: 150.0),
            payNowButton.rightAnchor.constraint(equalTo: view.centerXAnchor, constant: -9.0),
            payNowButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
            cancelButton.leftAnchor.constraint(equalTo: view.centerXAnchor, constant: 9.0),
            cancelButton.widthAnchor.constraint(equalToConstant: 150.0),
            cancelButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])
    }

    func paymentCardTextFieldDidChange(_ textField: STPPaymentCardTextField) {
        if cardEntry.isValid {
            payNowButton.isEnabled = true
            payNowButton.backgroundColor = scholaBlue
        } else {
            payNowButton.isEnabled = false
            payNowButton.backgroundColor = scholaBlueDisabled
        }
    }

    func paymentCardTextFieldWillEndEditing(forReturn textField: STPPaymentCardTextField) {
        if cardEntry.isValid {
            payNowButton(payNowButton)
        }
    }

    @objc func payNowButton(_ sender: UIButton) {
        payNowButton.setTitle("Paying...", for: .normal)
        payNowButton.isEnabled = false
        payNowButton.backgroundColor = scholaBlueDisabled
        let params = STPPaymentMethodParams(card: STPPaymentMethodCardParams(cardSourceParams: cardEntry.cardParams), billingDetails: nil, metadata: nil)
        STPAPIClient.shared().createPaymentMethod(with: params) { method, error in
            DispatchQueue.main.async {
                if let error = error {
                    self.payNowButton.setTitle("Pay Now", for: .normal)
                    self.payNowButton.isEnabled = true
                    self.payNowButton.backgroundColor = scholaBlue
                    let alert = UIAlertController(title: "Payment Error", message: error.localizedDescription, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    return
                }
                self.order.payments[0].method = method!.stripeId
                backend.placeOrder(order: self.order) { order, error in
                    DispatchQueue.main.async {
                        if let error = error {
                            self.payNowButton.setTitle("Pay Now", for: .normal)
                            self.payNowButton.isEnabled = true
                            self.payNowButton.backgroundColor = scholaBlue
                            let alert = UIAlertController(title: "Payment Error", message: error, preferredStyle: .alert)
                            alert.addAction(UIAlertAction(title: "OK", style: .default))
                            self.present(alert, animated: true, completion: nil)
                            return
                        }
                        for line in self.order.lines {
                            store.admitted += line.used ?? 0
                            store.sold += line.quantity
                        }
                        if self.order.email ?? "" != "" {
                            self.navigationController!.popToRootViewController(animated: true)
                            return
                        }
                        self.navigationController!.pushViewController(SalesReceipt(order: order!), animated: true)
                    }
                }
            }
        }
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }
}
