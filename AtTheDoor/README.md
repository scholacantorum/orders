# AtTheDoor

AtTheDoor is the iOS application used by Schola Cantorum for serving patrons at
the door of a concert.  Its primary purpose is selling tickets, for cash, check,
or card, using the Stripe card reader (or manual entry).  It can also consume
tickets bought it advance, by scanning the ticket barcode, entering the order
number, or searching on the customer name.

AtTheDoor is a Vue Native application.  As such, it is written in JavaScript,
with the Vue programming model and the React Native component set, and compiled
into an iOS application using the Xcode toolset.  It uses three major libraries:
`react-native-camera` for scanning QR codes, `react-native-stripe-terminal` for
interacting with the Stripe card reader, and `tipsi-stripe` for manual-entry
card payments.  The reason that AtTheDoor is an iOS native application is
because that is the only environment, at this writing, that supports both a
JavaScript/desired-state programming environment (Vue or React) and also our
Stripe card reader (BBPOS Chipper 2X BT).

* *Why a native phone/tablet app rather than a webapp?* Because Stripe only
  supports our card reader from native apps.  The other reader they support
  (Verifone P400) can be used from webapps, but it's five times as expensive and
  it requires hard-wired power and Internet.
* *Why Vue Native?* JavaScript, and the Vue programming model, are used in all
  of Schola's other webapps, so impose less knowledge burden on anyone needing
  to maintain them all.  Also, Vue Native can nominally be used on both iOS and
  Android, which could be handy someday.
* *Why iOS only, and not also Android?* Right now `react-native-stripe-terminal`
  is iOS only.  That will probably change fairly soon.

## Installation and Build

These installation instructions are for MacOS.  Other operating systems can be
used, but you'll need to consult the documentation for the underlying frameworks
to figure out how.

### Install basic tools

In the Mac App Store, install Xcode.  In its Preferences dialog, on the
Locations tab, select the most recent version in the Command Line Tools
dropdown.  In its Preferences dialog, on the Accounts tab, if your Apple ID is
not already listed, click the + button and add it.

If you don't already have them, install HomeBrew, Node, Yarn, Watchman, the Java
Developer Kit, React Native CLI, Vue Native CLI, and CocoaPods, with the
respective commands below.  If you do already have them, consider upgrading them
to current versions.

```sh
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
$ brew install node
$ brew install yarn
$ brew install watchman
$ brew tap AdoptOpenJDK/openjdk && brew cask install adoptopenjdk8
$ yarn global add react-native-cli
$ yarn global add vue-native-cli
$ sudo gem install cocoapods
```

### Install build dependencies

Install the JavaScript dependencies by running `yarn install` in the `AtTheDoor`
directory.

Install the iOS native dependencies by running `pod install` in the
`AtTheDoor/ios` directory.  Note that this can take upwards of half an hour.

Start the Xcode app and open `AtTheDoor/ios/AtTheDoor.xcworkspace`.  In the
General tab, under Signing, change the Team to your Apple ID.

### Build and run the app

To run on an iPhone simulator on your Mac:

* Run `react-native run-ios` in the `AtTheDoor` directory.

To run on a real iOS device in debug mode:

* Start the Xcode app and open `AtTheDoor/ios/AtTheDoor.xcworkspace`.
* Connect the device with a USB cable, and when it asks, tell it that it can
  trust the computer.
* Choose Window > Devices and Simulators from the menu bar, and select your
  device from the list.  The first time you select it, you'll need to wait about
  half an hour while it generates debug information.
* Press the play button at the top left of the Xcode window.  The build may take as
  much as a half hour.
* The first time, you'll get an error that it couldn't run the app on your
  device because it needs permission.  Follow the instructions in that error
  message to add permissions, and then hit the play button again.

Note that the debug mode app will run only while connected to the Mac with the
USB cable.

To run on a real iOS device in production mode:

* Choose Product > Scheme > Edit Scheme from the menu.
* Change the build configuration to Release, and turn off the Debug executable
  checkbox.
* Then press play as above.  The build will take *much* longer.
* Once the app is installed on the device, you can disconnect the USB cable and
  it will continue to work.

## Maintenance Notes

`react-native-stripe-terminal`, as I found it, wasn't sufficient to our needs.
I had to update it for the current version of the Stripe Terminal SDK and also
add a bunch of functionality that we needed.  For that reason, this app is using
my fork of the original `react-native-stripe-terminal`.  I have a pull request
out to the author for it.  Once that's accepted, we should probably switch to
using the original rather than my fork.

Along similar lines, I get a build error from `react-native-stripe-terminal`
because of its import of `EventEmitter` in `connectionService.js`.  I have had
to change that manually to import `eventemitter3` instead.  This isn't even in
my fork because I'm waiting to hear back from the author on whether that's the
best solution to the problem.

## Project Creation

For future use in setting up similar projects, it's worth noting how this app
was created.  First, the build setup:

```sh
$ vue-native init AtTheDoor --no-crna
$ cd AtTheDoor
$ yarn add git+https://github.com/rothskeller/react-native-stripe-terminal.git
$ yarn add react-native-camera tipsi-stripe
$ cd ios
$ pod init
$ vi Podfile # set to current contents
$ pod install
```

Then I had to go to my Google Firebase dashboard, create the AtTheDoor project
there, and download `GoogleService-Info.plist` from it into the `ios` directory.

And then the Xcode project settings:

* Run Xcode on `AtTheDoor/ios/AtTheDoor.xcworkspace`
* Click on AtTheDoor (root in project explorer)
* On the General tab, change Bundle identifier to org.scholacantorum.orders.AtTheDoor.
* Under Signing, change Team to Steven Roth (Personal Team).
* Under Deployment Info, turn off the landscape options.
* On the far right, change organization to Schola Cantorum.
* Under AtTheDoor, select Info.plist.
* Change Privacy - Location When In Use Usage Description to "Location access is
  required in order to accept payments."
* Add "Required Background Modes", with one item under it set to "App
  comunnicates using CoreBluetooth".
* Add "Privacy - Bluetooth Peripheral Usage Description" set to "Bluetooth
  access is required to connect to the card reader."
* Add "Privacy - Microphone Usage Description" set to "Microphone access is
  required to connect to a card reader."
* Add "Privacy - Camera Usage Description" set to "Camera access is required to
  scan ticket barcodes."
* Save the Info.plist file.
* Right-click on the second `AtTheDoor` (first child of the root) and choose
  Add Files.  Add the `GoogleService-Info.plist` file.
* Click on the `AppDelegate.m` file.  Add `@import Firebase;` at the bottom of
  the imports and `[FIRApp configure];` as the first line of the method.  Save
  the file.
* Right-click on Libraries, choose Add Files to AtTheDoor
* Add node_modules/react-native-stripe-terminal/ios/RNStripeTerminal.xcodeproj
* Click on AtTheDoor (root in project explorer)
* On the Build Phases tab, under Link Binary With Libraries, add
  libRNStripeTerminal.a.
