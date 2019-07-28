//
//  UseOrder.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-27.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class UseOrder: UIViewController, UseTicketClassDelegate {

    var tokenOrID: String!
    var spinner: UIActivityIndicatorView!
    var usage: TicketUsage!
    var useButton: UIButton!

    init(tokenOrID: String) {
        super.init(nibName: nil, bundle: nil)
        self.tokenOrID = tokenOrID
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white
        spinner = UIActivityIndicatorView()
        spinner.style = .gray
        view.addSubview(spinner)
        NSLayoutConstraint.useAndActivateConstraints([
            spinner.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            spinner.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 18.0),
        ])
        spinner.startAnimating()
        backend.fetchTicketUsage(event: store.event, tokenOrID: tokenOrID) { usage, error in
            DispatchQueue.main.async {
                if let error = error {
                    let alert = UIAlertController(title: "Server Error", message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    self.navigationController!.popViewController(animated: true)
                    return
                }
                self.usage = usage!
                self.showUsage()
            }
        }
    }

    func showUsage() {
        spinner.removeFromSuperview()

        var constraints = [NSLayoutConstraint]()
        let orderBox = UIView()
        orderBox.backgroundColor = UIColor.lightGray
        view.addSubview(orderBox)
        constraints.append(contentsOf: [
            orderBox.leftAnchor.constraint(equalTo: view.leftAnchor),
            orderBox.rightAnchor.constraint(equalTo: view.rightAnchor),
            orderBox.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
        ])

        let orderIDLabel = UILabel()
        orderIDLabel.text = "Order #\(usage.id)"
        orderIDLabel.font = UIFont.systemFont(ofSize: 20.0)
        orderBox.addSubview(orderIDLabel)

        if let name = usage.name {
            let orderName = UILabel()
            orderName.text = name
            orderName.font = UIFont.boldSystemFont(ofSize: 24.0)
            orderBox.addSubview(orderName)
            constraints.append(contentsOf: [
                orderName.centerXAnchor.constraint(equalTo: orderBox.centerXAnchor),
                orderName.topAnchor.constraint(equalTo: orderBox.topAnchor, constant: 9.0),
                orderIDLabel.centerXAnchor.constraint(equalTo: orderBox.centerXAnchor),
                orderIDLabel.topAnchor.constraint(equalTo: orderName.bottomAnchor),
                orderIDLabel.bottomAnchor.constraint(equalTo: orderBox.bottomAnchor, constant: -9.0),
            ])
        } else {
            constraints.append(contentsOf: [
                orderIDLabel.centerXAnchor.constraint(equalTo: orderBox.centerXAnchor),
                orderIDLabel.topAnchor.constraint(equalTo: orderBox.topAnchor, constant: 9.0),
                orderIDLabel.bottomAnchor.constraint(equalTo: orderBox.bottomAnchor, constant: -9.0),
            ])
        }

        var topAnchor = orderBox.bottomAnchor
        for cl in usage.classes {
            let clview = UIView()
            let clcontroller = UseTicketClass(cl, self)
            clcontroller.view = clview
            clcontroller.viewDidLoad()
            addChild(clcontroller)
            view.addSubview(clview)
            constraints.append(contentsOf: [
                clview.leftAnchor.constraint(equalTo: view.leftAnchor),
                clview.rightAnchor.constraint(equalTo: view.rightAnchor),
                clview.topAnchor.constraint(equalTo: topAnchor),
            ])
            topAnchor = clview.bottomAnchor
        }

        useButton = UIButton()
        useButton.setTitle("Use Tickets", for: .normal)
        useButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        useButton.setTitleColor(UIColor.white, for: .normal)
        useButton.layer.cornerRadius = 5.0
        setEnableDisable()
        useButton.addTarget(self, action: #selector(useButton(_:)), for: .touchUpInside)
        view.addSubview(useButton)
        constraints.append(contentsOf: [
            useButton.widthAnchor.constraint(equalToConstant: 150.0),
            useButton.rightAnchor.constraint(equalTo: view.centerXAnchor, constant: -9.0),
            useButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
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
            cancelButton.widthAnchor.constraint(equalToConstant: 150.0),
            cancelButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])

        NSLayoutConstraint.useAndActivateConstraints(constraints)
    }

    @objc func useButton(_ sender: UIButton) {
        useButton.setTitle("Saving...", for: .normal)
        useButton.isEnabled = false
        useButton.backgroundColor = scholaBlueDisabled
        backend.useTickets(event: store.event, usage: usage) { error in
            DispatchQueue.main.async {
                if let error = error {
                    let alert = UIAlertController(title: "Server Error", message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    self.useButton.setTitle("Use Tickets", for: .normal)
                    self.setEnableDisable()
                    return
                }
                for cl in self.usage.classes {
                    store.admitted += cl.used - cl.min
                }
                self.navigationController?.popToRootViewController(animated: true)
            }
        }
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }

    func classUsedChange(_ ticketClass: String, _ used: Int) {
        for (index, cl) in usage.classes.enumerated() {
            if cl.name == ticketClass {
                usage.classes[index].used = used
            }
        }
        setEnableDisable()
    }

    func setEnableDisable() {
        var enabled = false
        for cl in usage.classes {
            if cl.used != cl.min {
                enabled = true
            }
        }
        if enabled {
            useButton.backgroundColor = scholaBlue
            useButton.isEnabled = true
        } else {
            useButton.backgroundColor = scholaBlueDisabled
            useButton.isEnabled = false
        }
    }

}
