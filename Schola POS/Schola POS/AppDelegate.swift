//
//  AppDelegate.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-19.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

    var window: UIWindow?

    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        window = UIWindow(frame: UIScreen.main.bounds)
        let root = Login()
        // let root = EventChooser()
        // let root = ChooseCardReader()
        // let root = Main()
        let navbar = TopNavigation(rootViewController: root)
        window!.rootViewController = navbar
        window!.makeKeyAndVisible()
        return true
    }

}

