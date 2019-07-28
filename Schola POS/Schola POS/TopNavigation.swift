//
//  NavigationController.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-24.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class TopNavigation: UINavigationController {

    override func viewDidLoad() {
        super.viewDidLoad()
        navigationBar.barTintColor = scholaBlue
        navigationBar.barStyle = UIBarStyle.black

        store.logoutCallback = {
            DispatchQueue.main.async {
                self.popToRootViewController(animated: true)
            }
        }
    }

}
