//
//  SalesReceipt.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-26.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class SalesReceipt: UIViewController {

    var order: Order!
    var emailTextField: UITextField!
    var sendReceiptButton: UIButton!

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

        let successLabel = UILabel()
        successLabel.text = "Order Successful"
        successLabel.textColor = scholaBlue
        successLabel.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.addSubview(successLabel)

        let receiptLabel = UILabel()
        receiptLabel.text = "Send email receipt?"
        receiptLabel.font = UIFont.systemFont(ofSize: 24.0)
        view.addSubview(receiptLabel)

        let emailLabel = UILabel()
        emailLabel.text = "Email address"
        emailLabel.setContentHuggingPriority(.required, for: .horizontal)
        view.addSubview(emailLabel)

        emailTextField = UITextField()
        emailTextField.text = order.email ?? ""
        emailTextField.textContentType = .emailAddress
        emailTextField.autocapitalizationType = .none
        emailTextField.autocorrectionType = .no
        emailTextField.keyboardType = .emailAddress
        emailTextField.enablesReturnKeyAutomatically = true
        emailTextField.returnKeyType = .done
        emailTextField.borderStyle = .roundedRect
        NotificationCenter.default.addObserver(self, selector: #selector(textChange(_:)), name: UITextField.textDidChangeNotification, object: emailTextField)
        view.addSubview(emailTextField)

        sendReceiptButton = UIButton()
        sendReceiptButton.setTitle("Send Receipt", for: .normal)
        sendReceiptButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        sendReceiptButton.setTitleColor(UIColor.white, for: .normal)
        sendReceiptButton.layer.cornerRadius = 5.0
        if emailPred.evaluate(with: order.email ?? "") {
            sendReceiptButton.backgroundColor = scholaBlue
        } else {
            sendReceiptButton.backgroundColor = scholaBlueDisabled
            sendReceiptButton.isEnabled = false
        }
        sendReceiptButton.addTarget(self, action: #selector(sendReceiptButton(_:)), for: .touchUpInside)
        view.addSubview(sendReceiptButton)

        let noReceiptButton = UIButton()
        noReceiptButton.setTitle("No Receipt", for: .normal)
        noReceiptButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        noReceiptButton.setTitleColor(UIColor.white, for: .normal)
        noReceiptButton.layer.cornerRadius = 5.0
        noReceiptButton.backgroundColor = UIColor.darkGray
        noReceiptButton.addTarget(self, action: #selector(noReceiptButton(_:)), for: .touchUpInside)
        view.addSubview(noReceiptButton)

        NSLayoutConstraint.useAndActivateConstraints([
            summaryView.leftAnchor.constraint(equalTo: view.leftAnchor),
            summaryView.rightAnchor.constraint(equalTo: view.rightAnchor),
            summaryView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 18.0),
            {
                let constraint = summaryView.heightAnchor.constraint(equalToConstant: 1.0)
                constraint.priority = .defaultLow
                return constraint
            }(),
            successLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            successLabel.topAnchor.constraint(equalTo: summaryView.bottomAnchor, constant: 36.0),
            receiptLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            receiptLabel.topAnchor.constraint(equalTo: successLabel.bottomAnchor, constant: 36.0),
            emailTextField.leftAnchor.constraint(equalTo: emailLabel.rightAnchor, constant: 9.0),
            emailTextField.rightAnchor.constraint(equalTo: view.rightAnchor, constant: -9.0),
            emailTextField.topAnchor.constraint(equalTo: receiptLabel.bottomAnchor, constant: 18.0),
            emailTextField.heightAnchor.constraint(equalToConstant: 31.0),
            emailLabel.leftAnchor.constraint(equalTo: view.leftAnchor, constant: 9.0),
            emailLabel.centerYAnchor.constraint(equalTo: emailTextField.centerYAnchor),
            sendReceiptButton.widthAnchor.constraint(equalToConstant: 150.0),
            sendReceiptButton.rightAnchor.constraint(equalTo: view.centerXAnchor, constant: -9.0),
            sendReceiptButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
            noReceiptButton.leftAnchor.constraint(equalTo: view.centerXAnchor, constant: 9.0),
            noReceiptButton.widthAnchor.constraint(equalToConstant: 150.0),
            noReceiptButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])
    }

    @objc func textChange(_ notification: NSNotification) {
        if emailPred.evaluate(with: emailTextField.text ?? "") {
            sendReceiptButton.backgroundColor = scholaBlue
            sendReceiptButton.isEnabled = true
        } else {
            sendReceiptButton.backgroundColor = scholaBlueDisabled
            sendReceiptButton.isEnabled = false
        }
    }

    @objc func sendReceiptButton(_ sender: UIButton) {
        sendReceiptButton.setTitle("Sending...", for: .normal)
        sendReceiptButton.isEnabled = false
        backend.sendEmailReceipt(order: order, email: emailTextField.text!) { error in
            DispatchQueue.main.async {
                if let error = error {
                    let alert = UIAlertController(title: "Server Error", message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    self.sendReceiptButton.isEnabled = true
                    self.sendReceiptButton.setTitle("Send Receipt", for: .normal)
                }
                self.navigationController!.popToRootViewController(animated: true)
            }
        }
    }

    @objc func noReceiptButton(_ sender: UIButton) {
        navigationController!.popToRootViewController(animated: true)
    }

    func validEmail(_ email: String) -> Bool {
        
        return false
    }
}
