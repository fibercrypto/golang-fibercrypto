import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12

RowLayout {
    spacing: 0

    property alias backButtonHide: toolButtonBack.hide
    property alias toolButtonBack: toolButtonBack
    property alias menuBarReal: menuBarReal
    property alias menuBarColor: menuBarReal.color

    property alias enableOutputs: menuItemOutputs.enabled
    property alias enablePendingTransactions: menuItemPendingTransactions.enabled
    property alias enableBlockchain: menuItemBlockchain.enabled
    property alias enableNetworking: menuItemNetworking.enabled
    property alias enableSettings: menuItemSettings.enabled

    // Signals
    signal outputsRequested()
    signal pendingTransactionsRequested()
    signal networkingRequested()
    signal settingsRequested()
    signal blockchainRequested()
    signal aboutRequested()
    signal aboutQtRequested()

    // Functions
    function back() {
        if (!toolButtonBack.hide) {
            toolButtonBack.clicked()
        }
    }

    ToolButton {
        id: toolButtonBack

        property bool hide: generalStackView.depth === 1 || (generalStackView.depth >= 1 && generalStackView.busy)

        focusPolicy: Qt.NoFocus

        // opacity
        opacity: hide ? 0 : 1
        Behavior on opacity { NumberAnimation { duration: 500; easing.type: Easing.OutQuint } }

        // position
        implicitHeight: menuBarReal.height
        implicitWidth: implicitHeight
        Layout.leftMargin: hide ? -(width + padding) : 0
        Behavior on Layout.leftMargin { NumberAnimation { duration: 500; easing.type: Easing.OutQuint } }
        z: 1

        // icon
        icon.source: "qrc:/images/resources/images/icons/back.svg"
        icon.color: "white"

        // PointingHandCursor
        MouseArea {
            anchors.fill: parent
            acceptedButtons: Qt.NoButton
            cursorShape: Qt.PointingHandCursor
        }

        onClicked: {
            generalStackView.pop()

            enableOutputs = true
            enablePendingTransactions = true
            enableBlockchain = true
            enableNetworking = true
            enableSettings = true
        }
    }

    MenuBar {
        id: menuBarReal

        readonly property color iconColor: "transparent"
        property color color: Material.dialogColor
        property color menuTextColor: toolButtonBack.hide ? Material.primaryTextColor : "white"

        Layout.fillWidth: true
        leftInset:  -(toolButtonBack.width + toolButtonBack.padding)
        Material.foreground: menuTextColor
        Behavior on menuTextColor { ColorAnimation { } }

        background: Rectangle {
            id: backgroundRectangle
            implicitHeight: 40
            color: toolButtonBack.hide ? Material.dialogColor : menuBarReal.color

            Behavior on color { ColorAnimation { } }
        }

        Menu {
            id: menuTools
            title: qsTr("&Tools")

            CustomMenuItem {
                id: menuItemOutputs
                text: qsTr("&Outputs")
                iconSource: "qrc:/images/resources/images/icons/outputs.svg"

                onClicked: outputsRequested()
            }
            CustomMenuItem {
                id: menuItemPendingTransactions
                text: qsTr("&Pending transactions")
                iconSource: "qrc:/images/resources/images/icons/pending-transactions.svg"

                onClicked: pendingTransactionsRequested()
            }
            CustomMenuItem {
                id: menuItemBlockchain
                text: qsTr("&Blockchain")
                iconSource: "qrc:/images/resources/images/icons/blockchain.svg"
                
                onClicked: blockchainRequested()
            }
            CustomMenuItem {
                id: menuItemNetworking
                text: qsTr("&Networking")
                iconSource: "qrc:/images/resources/images/icons/networking.svg"

                onClicked: networkingRequested()
            }
            CustomMenuItem {
                id: menuItemSettings
                text: qsTr("&Settings")
                iconSource: "qrc:/images/resources/images/icons/settings.svg"

                onClicked: settingsRequested()
            }
        } // menuTools
        Menu {
            id: menuHelp
            title: qsTr("&Help")
            CustomMenuItem {
                text: qsTr("&About FiberCrypto")
                iconSource: "qrc:/images/resources/images/icons/appIcon.png"

                onClicked: aboutRequested()
            }
            CustomMenuItem {
                text: qsTr("About &Qt")
                iconSource: "qrc:/images/resources/images/icons/qt_logo_green_rgb_256x256.png"

                onClicked: aboutQtRequested()
            }
        } // menuHelp
    } // menuBarReal
} // RowLayout (menuBar)
