//
//  SalesPaymentCash.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class SalesPaymentCash: UIViewController {

    var order: Order!
    var received = 0

    lazy var cashReceivedTextField: UITextField = {
        let view = UITextField()
        view.keyboardType = .numberPad
        view.enablesReturnKeyAutomatically = true
        view.returnKeyType = .done
        view.borderStyle = .roundedRect
        NotificationCenter.default.addObserver(self, selector: #selector(textChange(_:)), name: UITextField.textDidChangeNotification, object: view)
        return view
    }()
    lazy var changeAmountLabel: UILabel = {
        let view = UILabel()
        view.text = "$0"
        return view
    }()
    lazy var donateSwitch: UISwitch = {
        let view = UISwitch()
        return view
    }()
    lazy var confirmButton: UIButton = {
        let view = UIButton()
        view.setTitle("Confirm", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(confirmButton(_:)), for: .touchUpInside)
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

        var constraints: [NSLayoutConstraint] = []
        let summaryController = SalesSummary(order: order)
        let summaryView = UIView()
        summaryController.view = summaryView
        summaryController.viewDidLoad()
        addChild(summaryController)
        view.addSubview(summaryView)
        summaryController.didMove(toParent: self)
        constraints.append(contentsOf: [
            summaryView.leftAnchor.constraint(equalTo: view.leftAnchor),
            summaryView.rightAnchor.constraint(equalTo: view.rightAnchor),
            summaryView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 18.0),
        ])

        let vbox = UIView()
        view.addSubview(vbox)
        constraints.append(contentsOf: [
            vbox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            vbox.topAnchor.constraint(equalTo: summaryView.bottomAnchor, constant: 24.0),
        ])

        let amountDueLabel1 = UILabel()
        amountDueLabel1.text = "Amount Due"
        vbox.addSubview(amountDueLabel1)
        let amountDueLabel2 = UILabel()
        amountDueLabel2.text = "$\(order.payments[0].amount/100)"
        vbox.addSubview(amountDueLabel2)
        constraints.append(contentsOf: [
            amountDueLabel1.leftAnchor.constraint(equalTo: vbox.leftAnchor),
            amountDueLabel1.widthAnchor.constraint(equalToConstant: 150.0),
            amountDueLabel1.topAnchor.constraint(equalTo: vbox.topAnchor),
            amountDueLabel2.leftAnchor.constraint(equalTo: amountDueLabel1.rightAnchor, constant: 9.0),
            amountDueLabel2.rightAnchor.constraint(equalTo: vbox.rightAnchor),
            amountDueLabel2.topAnchor.constraint(equalTo: vbox.topAnchor),
        ])

        let cashReceivedHBox = UIView()
        vbox.addSubview(cashReceivedHBox)
        constraints.append(contentsOf: [
            cashReceivedHBox.leftAnchor.constraint(equalTo: vbox.leftAnchor),
            cashReceivedHBox.rightAnchor.constraint(equalTo: vbox.rightAnchor),
            cashReceivedHBox.topAnchor.constraint(equalTo: amountDueLabel1.bottomAnchor, constant: 18.0),
        ])

        let cashReceivedLabel = UILabel()
        cashReceivedLabel.text = "Cash Received"
        cashReceivedHBox.addSubview(cashReceivedLabel)
        constraints.append(contentsOf: [
            cashReceivedLabel.leftAnchor.constraint(equalTo: cashReceivedHBox.leftAnchor),
            cashReceivedLabel.widthAnchor.constraint(equalToConstant: 150.0),
            cashReceivedLabel.centerYAnchor.constraint(equalTo: cashReceivedHBox.centerYAnchor),
        ])

        cashReceivedHBox.addSubview(cashReceivedTextField)
        constraints.append(contentsOf: [
            cashReceivedTextField.leftAnchor.constraint(equalTo: cashReceivedLabel.rightAnchor, constant: 9.0),
            cashReceivedTextField.widthAnchor.constraint(equalToConstant: 50.0),
            cashReceivedTextField.centerYAnchor.constraint(equalTo: cashReceivedHBox.centerYAnchor),
            cashReceivedTextField.heightAnchor.constraint(equalToConstant: 31.0),
            cashReceivedHBox.heightAnchor.constraint(greaterThanOrEqualToConstant: 31.0), // in case no buttons
        ])
        var leftAnchor = cashReceivedTextField.rightAnchor

        let amount = order.payments[0].amount / 100
        received = amount
        cashReceivedTextField.text = "\(received)"
        let a5 = ((amount + 4) / 5) * 5 // round up to nearest 5
        if a5 != amount {
            let a5Button = UIButton()
            a5Button.setTitle("$\(a5)", for: .normal)
            a5Button.setTitleColor(UIColor.white, for: .normal)
            a5Button.layer.cornerRadius = 5.0
            a5Button.backgroundColor = scholaBlue
            a5Button.addTarget(self, action: #selector(cashAmountButton(_:)), for: .touchUpInside)
            cashReceivedHBox.addSubview(a5Button)
            constraints.append(contentsOf: [
                a5Button.leftAnchor.constraint(equalTo: leftAnchor, constant: 9.0),
                a5Button.topAnchor.constraint(equalTo: cashReceivedHBox.topAnchor),
                a5Button.bottomAnchor.constraint(equalTo: cashReceivedHBox.bottomAnchor),
                a5Button.widthAnchor.constraint(equalTo: a5Button.titleLabel!.widthAnchor, constant: 18.0),
            ])
            leftAnchor = a5Button.rightAnchor
        }
        let a20 = ((amount + 19) / 20) * 20 // round up to nearest 20
        if a5 != a20 {
            let a20Button = UIButton()
            a20Button.setTitle("$\(a20)", for: .normal)
            a20Button.setTitleColor(UIColor.white, for: .normal)
            a20Button.layer.cornerRadius = 5.0
            a20Button.backgroundColor = scholaBlue
            a20Button.addTarget(self, action: #selector(cashAmountButton(_:)), for: .touchUpInside)
            cashReceivedHBox.addSubview(a20Button)
            constraints.append(contentsOf: [
                a20Button.leftAnchor.constraint(equalTo: leftAnchor, constant: 9.0),
                a20Button.topAnchor.constraint(equalTo: cashReceivedHBox.topAnchor),
                a20Button.bottomAnchor.constraint(equalTo: cashReceivedHBox.bottomAnchor),
                a20Button.widthAnchor.constraint(equalTo: a20Button.titleLabel!.widthAnchor, constant: 18.0),
            ])
            leftAnchor = a20Button.rightAnchor
        }
        constraints.append(cashReceivedHBox.rightAnchor.constraint(greaterThanOrEqualTo: leftAnchor))

        let changeHBox = UIView()
        vbox.addSubview(changeHBox)
        constraints.append(contentsOf: [
            changeHBox.leftAnchor.constraint(equalTo: vbox.leftAnchor),
            changeHBox.rightAnchor.constraint(equalTo: vbox.rightAnchor),
            changeHBox.topAnchor.constraint(equalTo: cashReceivedHBox.bottomAnchor, constant: 18.0),
            changeHBox.bottomAnchor.constraint(equalTo: vbox.bottomAnchor),
        ])

        let changeLabel = UILabel()
        changeLabel.text = "Change"
        changeHBox.addSubview(changeLabel)
        constraints.append(contentsOf: [
            changeLabel.leftAnchor.constraint(equalTo: changeHBox.leftAnchor),
            changeLabel.widthAnchor.constraint(equalToConstant: 150.0),
            changeLabel.centerYAnchor.constraint(equalTo: changeHBox.centerYAnchor),
        ])

        changeHBox.addSubview(changeAmountLabel)
        constraints.append(contentsOf: [
            changeAmountLabel.leftAnchor.constraint(equalTo: changeLabel.rightAnchor, constant: 9.0),
            changeAmountLabel.widthAnchor.constraint(equalToConstant: 50.0),
            changeAmountLabel.centerYAnchor.constraint(equalTo: changeHBox.centerYAnchor),
        ])

        changeHBox.addSubview(donateSwitch)
        constraints.append(contentsOf: [
            donateSwitch.leftAnchor.constraint(equalTo: changeAmountLabel.rightAnchor, constant: 9.0),
            donateSwitch.topAnchor.constraint(equalTo: changeHBox.topAnchor),
            donateSwitch.bottomAnchor.constraint(equalTo: changeHBox.bottomAnchor),
        ])

        let donateLabel = UILabel()
        donateLabel.text = "Donation"
        changeHBox.addSubview(donateLabel)
        constraints.append(contentsOf: [
            donateLabel.leftAnchor.constraint(equalTo: donateSwitch.rightAnchor, constant: 9.0),
            donateLabel.centerYAnchor.constraint(equalTo: changeHBox.centerYAnchor),
            changeHBox.rightAnchor.constraint(greaterThanOrEqualTo: donateLabel.rightAnchor),
        ])

        view.addSubview(confirmButton)
        constraints.append(contentsOf: [
            confirmButton.widthAnchor.constraint(equalToConstant: 100.0),
            confirmButton.rightAnchor.constraint(equalTo: view.centerXAnchor, constant: -9.0),
            confirmButton.topAnchor.constraint(equalTo: vbox.bottomAnchor, constant: 18.0),
        ])

        let cancelButton = UIButton()
        cancelButton.setTitle("Cancel", for: .normal)
        cancelButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        cancelButton.setTitleColor(UIColor.white, for: .normal)
        cancelButton.layer.cornerRadius = 5.0
        cancelButton.backgroundColor = UIColor.darkGray
        cancelButton.addTarget(self, action: #selector(cancelButton(_:)), for: .touchUpInside)
        view.addSubview(cancelButton)
        constraints.append(contentsOf: [
            cancelButton.leftAnchor.constraint(equalTo: view.centerXAnchor, constant: 9.0),
            cancelButton.widthAnchor.constraint(equalToConstant: 100.0),
            cancelButton.topAnchor.constraint(equalTo: vbox.bottomAnchor, constant: 18.0),
        ])

        NSLayoutConstraint.useAndActivateConstraints(constraints)
    }
    
    @objc func textChange(_ notification: NSNotification) {
        received = Int(cashReceivedTextField.text ?? "") ?? 0
        setChange(received)
    }

    @objc func cashAmountButton(_ sender: UIButton) {
        let receivedTP = sender.currentTitle!
        let receivedT = receivedTP[receivedTP.index(after: receivedTP.startIndex)...]
        received = Int(receivedT) ?? 0
        cashReceivedTextField.text = "\(received)"
        setChange(received)
    }

    func setChange(_ received: Int) {
        let change = received - (order.payments[0].amount / 100)
        if change > 0 {
            changeAmountLabel.text = "$\(change)"
            donateSwitch.isEnabled = true
            confirmButton.isEnabled = true
            confirmButton.backgroundColor = scholaBlue
        } else if change == 0 {
            changeAmountLabel.text = "$0"
            donateSwitch.isEnabled = false
            confirmButton.isEnabled = true
            confirmButton.backgroundColor = scholaBlue
        } else {
            changeAmountLabel.text = ""
            donateSwitch.isEnabled = false
            confirmButton.isEnabled = false
            confirmButton.backgroundColor = scholaBlueDisabled
        }
    }

    @objc func confirmButton(_ sender: UIButton) {
        confirmButton.isEnabled = false
        confirmButton.setTitle("Saving...", for: .normal)
        var donation = 0
        var donationLine: OrderLine!
        if donateSwitch.isEnabled && donateSwitch.isOn {
            donation = received*100 - order.payments[0].amount
            for line in order.lines {
                if line.product == "donation" {
                    donationLine = line
                    break
                }
            }
            if donationLine == nil {
                donationLine = OrderLine(product: "donation", quantity: 1, used: nil, usedAt: nil, price: 0)
                order.lines.append(donationLine)
            }
            donationLine.price += donation
            order.payments[0].amount += donation
        }
        backend.placeOrder(order: order) { placed, error in
            DispatchQueue.main.async {
                if let error = error {
                    if donation != 0 { // remove donation from failed order
                        self.order.payments[0].amount -= donation
                        donationLine.price -= donation
                        if donationLine.price == 0 {
                            self.order.lines.removeLast()
                        }
                    }
                    let alert = UIAlertController(title: nil, message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    return
                }
                for line in self.order.lines {
                    store.admitted += line.used ?? 0
                    if line.product != "donation" {
                        store.sold += line.quantity
                    }
                }
                store.cash += self.order.payments[0].amount
                self.navigationController!.popToRootViewController(animated: true)
            }
        }
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }

}
