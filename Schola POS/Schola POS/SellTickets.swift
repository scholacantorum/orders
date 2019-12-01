//
//  SellTicketsController.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-24.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import StripeTerminal

let scholaBlueDisabled = UIColor(red: 84.0 / 255.0, green: 169.0 / 255.0, blue: 253.0 / 255.0, alpha: 1.0)
// 50% lighter than scholaBlue

let nameRE = ".*\\S.*"
let namePred = NSPredicate(format: "SELF MATCHES %@", nameRE)
let emailRE = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
let emailPred = NSPredicate(format: "SELF MATCHES %@", emailRE)

class SellTickets: UIViewController, TicketQuantityDelegate, UITextFieldDelegate {

    var sellqty: [Int] = []
    var useqty: [Int] = []
    var total = 0
    var showingNameEmail = false

    lazy var totalBorder: UIView = {
        let view = UIView(frame: .zero)
        view.backgroundColor = UIColor.black
        return view
    }()
    lazy var totalLabel: UILabel = {
        let view = UILabel()
        view.font = UIFont.systemFont(ofSize: 24.0)
        view.textAlignment = .right
        return view
    }()
    lazy var nameHBox = UIView()
    lazy var nameLabel: UILabel = {
        let view = UILabel()
        view.text = "Name"
        return view
    }()
    lazy var nameTextField: UITextField = {
        let view = UITextField()
        view.textContentType = .name
        view.autocapitalizationType = .words
        view.autocorrectionType = .no
        view.enablesReturnKeyAutomatically = true
        view.returnKeyType = .next
        view.borderStyle = .roundedRect
        NotificationCenter.default.addObserver(self, selector: #selector(textChange(_:)), name: UITextField.textDidChangeNotification, object: view)
        return view
    }()
    lazy var emailHBox = UIView()
    lazy var emailLabel: UILabel = {
        let view = UILabel()
        view.text = "Email"
        return view
    }()
    lazy var emailTextField: UITextField = {
        let view = UITextField()
        view.textContentType = .emailAddress
        view.autocapitalizationType = .none
        view.autocorrectionType = .no
        view.keyboardType = .emailAddress
        view.enablesReturnKeyAutomatically = true
        view.returnKeyType = .done
        view.borderStyle = .roundedRect
        NotificationCenter.default.addObserver(self, selector: #selector(textChange(_:)), name: UITextField.textDidChangeNotification, object: view)
        return view
    }()
    lazy var buttonBar = UIView()
    lazy var cashButton: UIButton = {
        let view = UIButton()
        view.setTitle("Cash", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlueDisabled
        view.isEnabled = false
        view.addTarget(self, action: #selector(cashButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var checkButton: UIButton = {
        let view = UIButton()
        view.setTitle("Check", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlueDisabled
        view.isEnabled = false
        view.addTarget(self, action: #selector(checkButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var cardButton: UIButton = {
        let view = UIButton()
        view.setTitle("Card", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlueDisabled
        view.isEnabled = false
        view.addTarget(self, action: #selector(cardButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var cancelButton: UIButton = {
        let view = UIButton()
        view.setTitle("Cancel", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = UIColor.darkGray
        view.addTarget(self, action: #selector(cancelButton(_:)), for: .touchUpInside)
        return view
    }()

    var noNameEmailConstraints: [NSLayoutConstraint] = []
    var nameEmailConstraints: [NSLayoutConstraint] = []

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white

        var topAnchor = view.safeAreaLayoutGuide.topAnchor
        var constraints: [NSLayoutConstraint] = []
        for product in store.products {
            sellqty.append(0)
            useqty.append(0)

            let qtyController = TicketQuantity(product: product, delegate: self)
            let qtyView = UIView()
            qtyController.view = qtyView
            qtyController.viewDidLoad()
            addChild(qtyController)
            view.addSubview(qtyView)
            qtyController.didMove(toParent: self)
            constraints.append(contentsOf: [
                qtyView.topAnchor.constraint(equalTo: topAnchor, constant: 18.0),
                qtyView.leftAnchor.constraint(equalTo: view.leftAnchor),
                qtyView.rightAnchor.constraint(equalTo: view.rightAnchor),
            ])
            topAnchor = qtyView.bottomAnchor
        }

        view.addSubview(totalBorder)
        view.addSubview(totalLabel)
        constraints.append(contentsOf: [
            totalBorder.topAnchor.constraint(equalTo: topAnchor, constant: 9.0),
            totalBorder.heightAnchor.constraint(equalToConstant: 1.0),
            totalBorder.leftAnchor.constraint(equalTo: totalLabel.leftAnchor),
            totalBorder.rightAnchor.constraint(equalTo: totalLabel.rightAnchor),
            totalLabel.topAnchor.constraint(equalTo: totalBorder.bottomAnchor),
            totalLabel.rightAnchor.constraint(equalTo: view.rightAnchor, constant: -9.0),
        ])

        view.addSubview(buttonBar)
        var leftAnchor = buttonBar.leftAnchor
        var leftOffset: CGFloat = 0.0
        if store.allow.cash {
            buttonBar.addSubview(cashButton)
            buttonBar.addSubview(checkButton)
            constraints.append(contentsOf: [
                cashButton.leftAnchor.constraint(equalTo: leftAnchor),
                cashButton.widthAnchor.constraint(equalToConstant: 90.0),
                cashButton.topAnchor.constraint(equalTo: buttonBar.topAnchor),
                checkButton.leftAnchor.constraint(equalTo: cashButton.rightAnchor, constant: 9.0),
                checkButton.widthAnchor.constraint(equalToConstant: 90.0),
                checkButton.topAnchor.constraint(equalTo: buttonBar.topAnchor),
            ])
            leftAnchor = checkButton.rightAnchor
            leftOffset = 9.0
        }
        if store.allow.card {
            buttonBar.addSubview(cardButton)
            constraints.append(contentsOf: [
                cardButton.leftAnchor.constraint(equalTo: leftAnchor, constant: leftOffset),
                cardButton.widthAnchor.constraint(equalToConstant: 90.0),
                cardButton.topAnchor.constraint(equalTo: buttonBar.topAnchor),
            ])
            leftAnchor = cardButton.rightAnchor
            leftOffset = 9.0
        }
        buttonBar.addSubview(cancelButton)
        constraints.append(contentsOf: [
            cancelButton.leftAnchor.constraint(equalTo: leftAnchor, constant: leftOffset),
            cancelButton.widthAnchor.constraint(equalToConstant: 90.0),
            cancelButton.rightAnchor.constraint(equalTo: buttonBar.rightAnchor),
            cancelButton.topAnchor.constraint(equalTo: buttonBar.topAnchor),
            cancelButton.bottomAnchor.constraint(equalTo: buttonBar.bottomAnchor),
            buttonBar.centerXAnchor.constraint(equalTo: view.centerXAnchor),
        ])

        NSLayoutConstraint.useAndActivateConstraints(constraints)
        noNameEmailConstraints = [
            buttonBar.topAnchor.constraint(equalTo: totalLabel.bottomAnchor, constant: 18.0),
        ]
        NSLayoutConstraint.useAndActivateConstraints(noNameEmailConstraints)
        totalLabel.text = "TOTAL   $0"
    }

    func addNameEmail() {
        if showingNameEmail {
            return
        }
        view.addSubview(nameHBox)
        nameHBox.addSubview(nameLabel)
        nameHBox.addSubview(nameTextField)
        view.addSubview(emailHBox)
        emailHBox.addSubview(emailLabel)
        emailHBox.addSubview(emailTextField)
        if nameEmailConstraints.count == 0 {
            nameEmailConstraints = [
                nameHBox.leftAnchor.constraint(equalTo: view.leftAnchor),
                nameHBox.rightAnchor.constraint(equalTo: view.rightAnchor),
                nameHBox.topAnchor.constraint(equalTo: totalLabel.bottomAnchor, constant: 18.0),
                nameLabel.leftAnchor.constraint(equalTo: nameHBox.leftAnchor, constant: 9.0),
                nameLabel.widthAnchor.constraint(equalToConstant: 45.0),
                nameLabel.centerYAnchor.constraint(equalTo: nameHBox.centerYAnchor),
                nameTextField.leftAnchor.constraint(equalTo: nameLabel.rightAnchor, constant: 9.0),
                nameTextField.rightAnchor.constraint(equalTo: nameHBox.rightAnchor, constant: -9.0),
                nameTextField.topAnchor.constraint(equalTo: nameHBox.topAnchor),
                nameTextField.bottomAnchor.constraint(equalTo: nameHBox.bottomAnchor),
                nameTextField.heightAnchor.constraint(equalToConstant: 31.0),
                emailHBox.leftAnchor.constraint(equalTo: view.leftAnchor),
                emailHBox.rightAnchor.constraint(equalTo: view.rightAnchor),
                emailHBox.topAnchor.constraint(equalTo: nameHBox.bottomAnchor, constant: 9.0),
                emailLabel.leftAnchor.constraint(equalTo: emailHBox.leftAnchor, constant: 9.0),
                emailLabel.widthAnchor.constraint(equalToConstant: 45.0),
                emailLabel.centerYAnchor.constraint(equalTo: emailHBox.centerYAnchor),
                emailTextField.leftAnchor.constraint(equalTo: emailLabel.rightAnchor, constant: 9.0),
                emailTextField.rightAnchor.constraint(equalTo: emailHBox.rightAnchor, constant: -9.0),
                emailTextField.topAnchor.constraint(equalTo: emailHBox.topAnchor),
                emailTextField.bottomAnchor.constraint(equalTo: emailHBox.bottomAnchor),
                emailTextField.heightAnchor.constraint(equalToConstant: 31.0),
                buttonBar.topAnchor.constraint(equalTo: emailHBox.bottomAnchor, constant: 18.0),
            ]
        }
        NSLayoutConstraint.deactivate(noNameEmailConstraints)
        NSLayoutConstraint.useAndActivateConstraints(nameEmailConstraints)
        showingNameEmail = true
    }

    func removeNameEmail() {
        if !showingNameEmail {
            return
        }
        NSLayoutConstraint.deactivate(nameEmailConstraints)
        NSLayoutConstraint.activate(noNameEmailConstraints)
        nameLabel.removeFromSuperview()
        nameTextField.removeFromSuperview()
        nameHBox.removeFromSuperview()
        emailLabel.removeFromSuperview()
        emailTextField.removeFromSuperview()
        emailHBox.removeFromSuperview()
        showingNameEmail = false
    }

    func ticketQuantityChange(product: Product, sellQty: Int, useQty: Int) {
        total = 0
        for (index, prod) in store.products.enumerated() {
            if prod.id == product.id {
                sellqty[index] = sellQty
                useqty[index] = useQty
            }
            total += prod.price * sellqty[index]
        }
        totalLabel.text = "TOTAL   $\(total/100)"
        enableDisable()
    }

    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        textField.resignFirstResponder()
        return true
    }

    func textFieldDidEndEditing(_ textField: UITextField) {
        if textField == nameTextField {
            emailTextField.becomeFirstResponder()
        }
        enableDisable()
    }

    @objc func textChange(_ notification: NSNotification) {
        enableDisable()
    }

    func enableDisable() {
        var showNameEmail = false
        var enableCash = false
        var enableCheckCard = false
        for (index, prod) in store.products.enumerated() {
            if useqty[index] < prod.ticketCount*sellqty[index] {
                showNameEmail = true
            }
            if sellqty[index] > 0 {
                enableCash = true
                if prod.price > 0 {
                    enableCheckCard = true
                }
            }
        }
        if showNameEmail {
            addNameEmail()
            if !namePred.evaluate(with: nameTextField.text ?? "") {
                enableCash = false
                enableCheckCard = false
            }
            if !emailPred.evaluate(with: emailTextField.text ?? "") {
                enableCash = false
                enableCheckCard = false
            }
        } else {
            removeNameEmail()
        }
        if store.allow.cash {
            cashButton.isEnabled = enableCash
            cashButton.backgroundColor = enableCash ? scholaBlue : scholaBlueDisabled
            checkButton.isEnabled = enableCheckCard
            checkButton.backgroundColor = enableCheckCard ? scholaBlue : scholaBlueDisabled
        }
        if store.allow.card {
            cardButton.isEnabled = enableCheckCard
            cardButton.backgroundColor = enableCheckCard ? scholaBlue : scholaBlueDisabled
        }
    }

    @objc func cashButton(_ sender: UIButton) {
        let order = createOrder(paymentType: "cash", paymentSubtype: "")
        navigationController!.pushViewController(SalesPaymentCash(order: order), animated: true)
    }

    @objc func checkButton(_ sender: UIButton) {
        let order = createOrder(paymentType: "check", paymentSubtype: "")
        navigationController!.pushViewController(SalesPaymentCheck(order: order), animated: true)
    }

    @objc func cardButton(_ sender: UIButton) {
        if Terminal.shared.connectionStatus == .connected {
            let order = createOrder(paymentType: "card-present", paymentSubtype: nil)
            navigationController!.pushViewController(SalesPaymentCardPresent(order: order), animated: true)
        } else {
            let order = createOrder(paymentType: "card", paymentSubtype: "manual")
            navigationController!.pushViewController(SalesPaymentCard(order: order), animated: true)
        }
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }

    func createOrder(paymentType: String, paymentSubtype: String?) -> Order {
        let name = showingNameEmail ? nameTextField.text : nil
        let email = showingNameEmail ? emailTextField.text : nil
        var payment = OrderPayment(type: paymentType, subtype: paymentSubtype, method: nil, amount: 0)
        var lines: [OrderLine] = []
        for (index, prod) in store.products.enumerated() {
            if sellqty[index] > 0 {
                lines.append(OrderLine(product: prod.id, quantity: sellqty[index], used: useqty[index], usedAt: store.event.id, price: prod.price))
                payment.amount += sellqty[index] * prod.price
            }
        }
        return Order(id: nil, source: "inperson", name: name, email: email, payments: [payment], lines: lines, error: nil)
    }
}
