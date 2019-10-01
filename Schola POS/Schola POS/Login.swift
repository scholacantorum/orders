//
//  LoginController.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-19.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

let scholaBlue = UIColor(red: 1.0 / 255.0, green: 83.0 / 255.0, blue: 165.0 / 255.0, alpha: 1.0)

class Login: UIViewController, UITextFieldDelegate {

    //MARK: Views

    lazy var pleaseLogInLabel: UILabel = {
        let view = UILabel()
        view.text = "Please log in."
        view.font = UIFont.boldSystemFont(ofSize: 24.0)
        return view
    }()
    lazy var usernameHBox = UIView()
    lazy var usernameLabel: UILabel = {
        let view = UILabel()
        view.text = "Username"
        return view
    }()
    lazy var usernameTextField: UITextField = {
        let view = UITextField()
        view.textContentType = .username
        view.autocapitalizationType = .none
        view.autocorrectionType = .no
        view.enablesReturnKeyAutomatically = true
        view.returnKeyType = .next
        view.borderStyle = .roundedRect
        return view
    }()
    lazy var passwordHBox = UIView()
    lazy var passwordLabel: UILabel = {
        let view = UILabel()
        view.text = "Password"
        return view
    }()
    lazy var passwordTextField: UITextField = {
        let view = UITextField()
        view.textContentType = .password
        view.autocapitalizationType = .none
        view.autocorrectionType = .no
        view.enablesReturnKeyAutomatically = true
        view.returnKeyType = .go
        view.isSecureTextEntry = true
        view.borderStyle = .roundedRect
        return view
    }()
    lazy var testModeHBox = UIView()
    lazy var testModeSwitch: UISwitch = {
        let view = UISwitch()
        return view
    }()
    lazy var testModeLabel: UILabel = {
        let view = UILabel()
        view.text = "Test mode"
        return view
    }()
    lazy var allowCardHBox = UIView()
    lazy var allowCardSwitch: UISwitch = {
        let view = UISwitch()
        return view
    }()
    lazy var allowCardLabel: UILabel = {
        let view = UILabel()
        view.text = "Allow credit/debit card sales"
        return view
    }()
    lazy var allowCashHBox = UIView()
    lazy var allowCashSwitch: UISwitch = {
        let view = UISwitch()
        return view
    }()
    lazy var allowCashLabel: UILabel = {
        let view = UILabel()
        view.text = "Allow cash/check sales"
        return view
    }()
    lazy var allowWillCallHBox = UIView()
    lazy var allowWillCallSwitch: UISwitch = {
        let view = UISwitch()
        return view
    }()
    lazy var allowWillCallLabel: UILabel = {
        let view = UILabel()
        view.text = "Allow will call"
        return view
    }()
    lazy var loginButton: UIButton = {
        let view = UIButton(type: .system)
        view.setTitle("Login", for: .normal)
        view.titleLabel!.font = UIFont(name: "System", size: 20.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(loginButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var errorLabel: UILabel = {
        let view = UILabel()
        view.text = ""
        view.textColor = UIColor.red
        return view
    }()

    var loggingIn = false

    override func viewDidLoad() {
        super.viewDidLoad()
        navigationItem.title = "Schola Point of Sale"
        view.backgroundColor = UIColor.white
        usernameTextField.delegate = self
        passwordTextField.delegate = self
        view.addSubview(pleaseLogInLabel)
        view.addSubview(usernameHBox)
        usernameHBox.addSubview(usernameLabel)
        usernameHBox.addSubview(usernameTextField)
        view.addSubview(passwordHBox)
        passwordHBox.addSubview(passwordLabel)
        passwordHBox.addSubview(passwordTextField)
        view.addSubview(testModeHBox)
        testModeHBox.addSubview(testModeSwitch)
        testModeHBox.addSubview(testModeLabel)
        view.addSubview(allowCardHBox)
        allowCardHBox.addSubview(allowCardSwitch)
        allowCardHBox.addSubview(allowCardLabel)
        view.addSubview(allowCashHBox)
        allowCashHBox.addSubview(allowCashSwitch)
        allowCashHBox.addSubview(allowCashLabel)
        view.addSubview(allowWillCallHBox)
        allowWillCallHBox.addSubview(allowWillCallSwitch)
        allowWillCallHBox.addSubview(allowWillCallLabel)
        view.addSubview(loginButton)
        view.addSubview(errorLabel)
        NSLayoutConstraint.useAndActivateConstraints([
            pleaseLogInLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            pleaseLogInLabel.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 18.0),
            usernameHBox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            usernameHBox.topAnchor.constraint(equalTo: pleaseLogInLabel.bottomAnchor, constant: 18.0),
            usernameLabel.leftAnchor.constraint(equalTo: usernameHBox.leftAnchor),
            usernameLabel.centerYAnchor.constraint(equalTo: usernameHBox.centerYAnchor),
            usernameTextField.leftAnchor.constraint(greaterThanOrEqualTo: usernameLabel.rightAnchor, constant: 18.0),
            usernameTextField.rightAnchor.constraint(equalTo: usernameHBox.rightAnchor),
            usernameTextField.topAnchor.constraint(equalTo: usernameHBox.topAnchor),
            usernameTextField.bottomAnchor.constraint(equalTo: usernameHBox.bottomAnchor),
            usernameTextField.widthAnchor.constraint(equalToConstant: 150.0),
            usernameTextField.heightAnchor.constraint(equalToConstant: 31.0),
            passwordHBox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            passwordHBox.topAnchor.constraint(equalTo: usernameHBox.bottomAnchor, constant: 9.0),
            passwordLabel.leftAnchor.constraint(equalTo: passwordHBox.leftAnchor),
            passwordLabel.centerYAnchor.constraint(equalTo: passwordHBox.centerYAnchor),
            passwordTextField.leftAnchor.constraint(greaterThanOrEqualTo: passwordLabel.rightAnchor, constant: 18.0),
            passwordTextField.rightAnchor.constraint(equalTo: passwordHBox.rightAnchor),
            passwordTextField.topAnchor.constraint(equalTo: passwordHBox.topAnchor),
            passwordTextField.bottomAnchor.constraint(equalTo: passwordHBox.bottomAnchor),
            passwordTextField.widthAnchor.constraint(equalToConstant: 150.0),
            passwordTextField.heightAnchor.constraint(equalToConstant: 31.0),
            passwordHBox.leftAnchor.constraint(equalTo: usernameHBox.leftAnchor),
            passwordHBox.rightAnchor.constraint(equalTo: usernameHBox.rightAnchor),
            testModeHBox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            testModeHBox.topAnchor.constraint(equalTo: passwordHBox.bottomAnchor, constant: 18.0),
            testModeSwitch.leftAnchor.constraint(equalTo: testModeHBox.leftAnchor),
            testModeSwitch.topAnchor.constraint(equalTo: testModeHBox.topAnchor),
            testModeSwitch.bottomAnchor.constraint(equalTo: testModeHBox.bottomAnchor),
            testModeLabel.leftAnchor.constraint(equalTo: testModeSwitch.rightAnchor, constant: 9.0),
            testModeLabel.rightAnchor.constraint(lessThanOrEqualTo: testModeHBox.rightAnchor),
            testModeLabel.centerYAnchor.constraint(equalTo: testModeHBox.centerYAnchor),
            allowCardHBox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            allowCardHBox.topAnchor.constraint(equalTo: testModeHBox.bottomAnchor, constant: 18.0),
            allowCardSwitch.leftAnchor.constraint(equalTo: allowCardHBox.leftAnchor),
            allowCardSwitch.topAnchor.constraint(equalTo: allowCardHBox.topAnchor),
            allowCardSwitch.bottomAnchor.constraint(equalTo: allowCardHBox.bottomAnchor),
            allowCardLabel.leftAnchor.constraint(equalTo: allowCardSwitch.rightAnchor, constant: 9.0),
            allowCardLabel.rightAnchor.constraint(lessThanOrEqualTo: allowCardHBox.rightAnchor),
            allowCardLabel.centerYAnchor.constraint(equalTo: allowCardHBox.centerYAnchor),
            allowCardHBox.leftAnchor.constraint(equalTo: testModeHBox.leftAnchor),
            allowCardHBox.rightAnchor.constraint(equalTo: testModeHBox.rightAnchor),
            allowCashHBox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            allowCashHBox.topAnchor.constraint(equalTo: allowCardHBox.bottomAnchor, constant: 9.0),
            allowCashSwitch.leftAnchor.constraint(equalTo: allowCashHBox.leftAnchor),
            allowCashSwitch.topAnchor.constraint(equalTo: allowCashHBox.topAnchor),
            allowCashSwitch.bottomAnchor.constraint(equalTo: allowCashHBox.bottomAnchor),
            allowCashLabel.leftAnchor.constraint(equalTo: allowCashSwitch.rightAnchor, constant: 9.0),
            allowCashLabel.rightAnchor.constraint(lessThanOrEqualTo: allowCashHBox.rightAnchor),
            allowCashLabel.centerYAnchor.constraint(equalTo: allowCashHBox.centerYAnchor),
            allowCashHBox.leftAnchor.constraint(equalTo: testModeHBox.leftAnchor),
            allowCashHBox.rightAnchor.constraint(equalTo: testModeHBox.rightAnchor),
            allowWillCallHBox.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            allowWillCallHBox.topAnchor.constraint(equalTo: allowCashHBox.bottomAnchor, constant: 9.0),
            allowWillCallSwitch.leftAnchor.constraint(equalTo: allowWillCallHBox.leftAnchor),
            allowWillCallSwitch.topAnchor.constraint(equalTo: allowWillCallHBox.topAnchor),
            allowWillCallSwitch.bottomAnchor.constraint(equalTo: allowWillCallHBox.bottomAnchor),
            allowWillCallLabel.leftAnchor.constraint(equalTo: allowWillCallSwitch.rightAnchor, constant: 9.0),
            allowWillCallLabel.rightAnchor.constraint(lessThanOrEqualTo: allowWillCallHBox.rightAnchor),
            allowWillCallLabel.centerYAnchor.constraint(equalTo: allowWillCallHBox.centerYAnchor),
            allowWillCallHBox.leftAnchor.constraint(equalTo: testModeHBox.leftAnchor),
            allowWillCallHBox.rightAnchor.constraint(equalTo: testModeHBox.rightAnchor),
            loginButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            loginButton.topAnchor.constraint(equalTo: allowWillCallHBox.bottomAnchor, constant: 18.0),
            loginButton.widthAnchor.constraint(equalToConstant: 150.0),
            loginButton.heightAnchor.constraint(equalToConstant: 30.0),
            errorLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            errorLabel.topAnchor.constraint(equalTo: loginButton.bottomAnchor, constant: 18.0),
        ])
    }

    override func viewWillAppear(_ animated: Bool) {
        usernameTextField.text = nil
        passwordTextField.text = nil
        testModeSwitch.isOn = false
        allowCardSwitch.isOn = true
        allowCashSwitch.isOn = true
        allowWillCallSwitch.isOn = true
    }

    //MARK: UITextFieldDelegate

    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        textField.resignFirstResponder()
        return true
    }

    func textFieldDidEndEditing(_ textField: UITextField) {
        if textField == usernameTextField {
            passwordTextField.becomeFirstResponder()
        } else {
            doLogin()
        }
    }

    //MARK: Actions

    @objc func loginButton(_ sender: UIButton) {
        doLogin()
    }

    func doLogin() {
        if loggingIn || store.auth != "" {
            return
        }
        let username = usernameTextField.text ?? ""
        let password = passwordTextField.text ?? ""
        if username == "" || password == "" {
            errorLabel.text = "Please enter username and password."
            return
        }
        let testmode = testModeSwitch.isOn
        var allow = Allow()
        allow.card = allowCardSwitch.isOn
        allow.cash = allowCashSwitch.isOn
        allow.willcall = allowWillCallSwitch.isOn
        loggingIn = true
        loginButton.isEnabled = false
        loginButton.setTitle("Logging in...", for: .normal)
        backend.login(username, password, testmode, allow) { error in
            self.loggingIn = false
            DispatchQueue.main.async {
                self.loginButton.isEnabled = true
                self.loginButton.setTitle("Login", for: .normal)
                if error != nil {
                    self.errorLabel.text = error
                } else {
                    self.navigationController?.pushViewController(EventChooser(), animated: true)
                }
            }
        }
    }
}

extension NSLayoutConstraint {

    public class func useAndActivateConstraints(_ constraints: [NSLayoutConstraint]) {
        for constraint in constraints {
            if let view = constraint.firstItem as? UIView {
                view.translatesAutoresizingMaskIntoConstraints = false
            }
        }
        activate(constraints)
    }
}
