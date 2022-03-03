//
//  SalesPaymentCheck.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class SalesPaymentCheck: UIViewController {

    var order: Order!
    var received = 0

    lazy var checkReceivedTextField: UITextField = {
        let view = UITextField()
        view.keyboardType = .numberPad
        view.enablesReturnKeyAutomatically = true
        view.returnKeyType = .done
        view.borderStyle = .roundedRect
        NotificationCenter.default.addObserver(self, selector: #selector(textChange(_:)), name: UITextField.textDidChangeNotification, object: view)
        return view
    }()
    lazy var donationAmountLabel: UILabel = {
        let view = UILabel()
        view.text = "$0"
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

        let checkReceivedHBox = UIView()
        vbox.addSubview(checkReceivedHBox)
        constraints.append(contentsOf: [
            checkReceivedHBox.leftAnchor.constraint(equalTo: vbox.leftAnchor),
            checkReceivedHBox.rightAnchor.constraint(equalTo: vbox.rightAnchor),
            checkReceivedHBox.topAnchor.constraint(equalTo: amountDueLabel1.bottomAnchor, constant: 18.0),
        ])

        let checkReceivedLabel = UILabel()
        checkReceivedLabel.text = "Check Amount"
        checkReceivedHBox.addSubview(checkReceivedLabel)
        constraints.append(contentsOf: [
            checkReceivedLabel.leftAnchor.constraint(equalTo: checkReceivedHBox.leftAnchor),
            checkReceivedLabel.widthAnchor.constraint(equalToConstant: 150.0),
            checkReceivedLabel.centerYAnchor.constraint(equalTo: checkReceivedHBox.centerYAnchor),
        ])

        checkReceivedHBox.addSubview(checkReceivedTextField)
        constraints.append(contentsOf: [
            checkReceivedTextField.leftAnchor.constraint(equalTo: checkReceivedLabel.rightAnchor, constant: 9.0),
            checkReceivedTextField.widthAnchor.constraint(equalToConstant: 50.0),
            checkReceivedTextField.centerYAnchor.constraint(equalTo: checkReceivedHBox.centerYAnchor),
            checkReceivedTextField.heightAnchor.constraint(equalToConstant: 31.0),
            checkReceivedHBox.heightAnchor.constraint(equalToConstant: 31.0),
            checkReceivedHBox.rightAnchor.constraint(greaterThanOrEqualTo: checkReceivedTextField.rightAnchor),
        ])

        let amount = order.payments[0].amount / 100
        received = amount
        checkReceivedTextField.text = "\(received)"

        let donationLabel = UILabel()
        donationLabel.text = "Donation"
        vbox.addSubview(donationLabel)
        constraints.append(contentsOf: [
            donationLabel.leftAnchor.constraint(equalTo: vbox.leftAnchor),
            donationLabel.widthAnchor.constraint(equalToConstant: 150.0),
            donationLabel.topAnchor.constraint(equalTo: checkReceivedHBox.bottomAnchor, constant: 18.0),
            donationLabel.bottomAnchor.constraint(equalTo: vbox.bottomAnchor),
        ])

        vbox.addSubview(donationAmountLabel)
        constraints.append(contentsOf: [
            donationAmountLabel.leftAnchor.constraint(equalTo: donationLabel.rightAnchor, constant: 9.0),
            donationAmountLabel.widthAnchor.constraint(equalToConstant: 50.0),
            donationAmountLabel.topAnchor.constraint(equalTo: checkReceivedHBox.bottomAnchor, constant: 18.0),
            donationAmountLabel.bottomAnchor.constraint(equalTo: vbox.bottomAnchor),
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
        received = Int(checkReceivedTextField.text ?? "") ?? 0
        setDonation(received)
    }

    @objc func checkAmountButton(_ sender: UIButton) {
        let receivedTP = sender.currentTitle!
        let receivedT = receivedTP[receivedTP.index(after: receivedTP.startIndex)...]
        received = Int(receivedT) ?? 0
        checkReceivedTextField.text = "\(received)"
        setDonation(received)
    }

    func setDonation(_ received: Int) {
        let donation = received - (order.payments[0].amount / 100)
        if donation >= 0 {
            donationAmountLabel.text = "$\(donation)"
            confirmButton.isEnabled = true
            confirmButton.backgroundColor = scholaBlue
        } else {
            donationAmountLabel.text = ""
            confirmButton.isEnabled = false
            confirmButton.backgroundColor = scholaBlueDisabled
        }
    }

    @objc func confirmButton(_ sender: UIButton) {
        confirmButton.isEnabled = false
        confirmButton.setTitle("Saving...", for: .normal)
        let donation = received*100 - order.payments[0].amount
        var donationLine: OrderLine!
        for line in order.lines {
            if line.product == "donation" {
                donationLine = line
                break
            }
        }
        if donation > 0 {
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
                    if donation >= 0 { // remove donation from failed order
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
                store.check += self.order.payments[0].amount
                self.navigationController?.popToRootViewController(animated: true)
            }
        }
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }

}
