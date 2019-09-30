//
//  UseTicketClass.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-27.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

let restrictedClassColor = UIColor(red: 1.0, green: 215.0 / 255.0, blue: 0.0, alpha: 1.0)
let goodTicketColor = UIColor(red: 40.0 / 255.0, green: 167.0 / 255.0, blue: 69 / 255.0, alpha: 1.0)

protocol UseTicketClassDelegate {
    func classUsedChange(_ ticketClass: String, _ used: Int)
}

class UseTicketClass: UIViewController {

    var usage: TicketClassUsage!
    var delegate: UseTicketClassDelegate!
    var buttons = [UIButton]()

    init(_ usage: TicketClassUsage, _ delegate: UseTicketClassDelegate) {
        super.init(nibName: nil, bundle: nil)
        self.usage = usage
        self.delegate = delegate
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        if usage.overflow ?? false {
            view.backgroundColor = UIColor.red
        } else if usage.name != "" {
            view.backgroundColor = restrictedClassColor
        } else {
            view.backgroundColor = UIColor.white
        }

        var constraints = [NSLayoutConstraint]()

        let inner = UIView()
        view.addSubview(inner)
        constraints.append(contentsOf: [
            inner.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            inner.widthAnchor.constraint(equalToConstant: 315.0),
            inner.topAnchor.constraint(equalTo: view.topAnchor, constant: 9.0),
            inner.bottomAnchor.constraint(equalTo: view.bottomAnchor, constant: -9.0),
        ])

        let nameLabel = UILabel()
        nameLabel.text = usage.name != "" ? usage.name : "General Admission"
        inner.addSubview(nameLabel)
        constraints.append(contentsOf: [
            nameLabel.leftAnchor.constraint(equalTo: inner.leftAnchor),
            nameLabel.topAnchor.constraint(equalTo: inner.topAnchor),
        ])

        var topAnchor = nameLabel.bottomAnchor
        var topOffset: CGFloat = 0.0
        var leftAnchor = inner.leftAnchor
        var leftOffset: CGFloat = 0.0
        let max = usage.max < 1000 ? usage.max : ((usage.used + 6) / 6) * 6
        for bnum in 1...max {
            let button = UIButton()
            if bnum <= usage.min {
                button.setTitle("Used", for: .normal)
                button.titleLabel!.font = UIFont.boldSystemFont(ofSize: 16.0)
            } else {
                button.setTitle("\(bnum - usage.min)", for: .normal)
                button.titleLabel!.font = UIFont.boldSystemFont(ofSize: 24.0)
            }
            button.layer.cornerRadius = 5.0
            button.addTarget(self, action: #selector(usageChange(_:)), for: .touchUpInside)
            inner.addSubview(button)
            constraints.append(contentsOf: [
                button.leftAnchor.constraint(equalTo: leftAnchor, constant: leftOffset),
                button.widthAnchor.constraint(equalToConstant: 45.0),
                button.topAnchor.constraint(equalTo: topAnchor, constant: topOffset),
                button.heightAnchor.constraint(equalToConstant: 45.0),
            ])
            if bnum % 6 == 0 {
                topAnchor = button.bottomAnchor
                topOffset = 9.0
                leftAnchor = inner.leftAnchor
                leftOffset = 0.0
            } else {
                leftAnchor = button.rightAnchor
                leftOffset = 9.0
            }
            buttons.append(button)
        }
        constraints.append(inner.bottomAnchor.constraint(equalTo: buttons.last!.bottomAnchor))

        NSLayoutConstraint.useAndActivateConstraints(constraints)
        setButtons()
    }

    @objc func usageChange(_ sender: UIButton) {
        let bnum = sender.currentTitle == "Used" ? usage.min : usage.min + Int(sender.currentTitle!)!
        if bnum < usage.min {
            return
        }
        if usage.overflow ?? false {
            usage.overflow = false
        }
        if bnum == 1 && usage.min == 0 && usage.used == 1 {
            usage.used = 0
        } else {
            usage.used = bnum
        }
        setButtons()
        delegate.classUsedChange(usage.name, usage.used)
    }

    func setButtons() {
        for (index, button) in buttons.enumerated() {
            let bnum = index+1
            if bnum <= usage.min {
                button.setTitleColor(UIColor.white, for: .normal)
                button.backgroundColor = UIColor.darkGray
                button.isEnabled = bnum == usage.min
            } else if bnum <= usage.used {
                button.setTitleColor(UIColor.white, for: .normal)
                button.backgroundColor = goodTicketColor
            } else {
                button.setTitleColor(goodTicketColor, for: .normal)
                button.backgroundColor = UIColor.white
                button.layer.borderColor = goodTicketColor.cgColor
                button.layer.borderWidth = 2.0
            }
        }
    }

}
