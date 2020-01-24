import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12

Menu {
    id: menuThemeAccent

    readonly property int currentTheme: applicationWindow.Material.theme
    readonly property color currentAccent: applicationWindow.accentColor
    property int currentSelectedIndex // initialized when the component is completed (see bellow)
    readonly property var materialPredefinedColors: [
        Material.Pink,
        Material.Purple,
        Material.DeepPurple,
        Material.Indigo,
        Material.Blue,
        Material.LightBlue,
        Material.Cyan,
        Material.Teal,
        Material.Green,
        Material.LightGreen,
        Material.Lime,
        Material.Yellow,
        Material.Amber,
        Material.Orange,
        Material.DeepOrange,
        Material.Brown,
        Material.Grey,
        Material.BlueGrey
    ]

    function saveConfiguration(theme, accent) {
        settings.setValue("style/material/theme", theme)
        settings.setValue("style/material/accent", accent)
    }

    // Initialize `currentSelectedIndex`:
    onAboutToShow: {
        for (var i = 0; i < gridAccents.children.length; i++) {
            if (gridAccents.children[i].checked) {
                currentSelectedIndex = i
            }
        }
    }

    title: qsTr("Style configuration")
    width: gridAccents.width
        
    SwitchDelegate {
        id: switchDelegateTheme

        width: parent.width
        text: qsTr("Night theme")
        icon.source: "qrc:/images/resources/images/icons/moon.svg"
        icon.color: "transparent"
        checked: currentTheme === Material.Dark

        onClicked: {
            applicationWindow.flash()
            applicationWindow.Material.theme = (currentTheme === Material.Light ? Material.Dark : Material.Light)
            applicationWindow.skipAccentColorAnimation = true
            applicationWindow.accentColor = gridAccents.children[currentSelectedIndex].Material.accent // needed for accent color transitions
            skipAccentColorAnimation = false
            saveConfiguration(Material.theme, Material.accentColor)
        }
    }

    MenuSeparator {}

    Grid {
        id: gridAccents

        property int markedIndex: 5

        columns: 4
        columnSpacing: 10
        rowSpacing: 10
        leftPadding: 10
        rightPadding: leftPadding
        bottomPadding: 4

        Repeater {
            model: 19

            delegate: Rectangle {
                id: rectangleDelegate

                property bool checked: applicationWindow.accentColor === Material.accent

                width: 48
                height: width
                radius: width/2
                color: Material.accent
                border { width: 3; color: Qt.darker(rectangleDelegate.color) }
                Material.accent: materialPredefinedColors[index]

                Rectangle {
                    id: rectangleCheckIndicator

                    anchors.centerIn: parent
                    width: rectangleDelegate.width - 16
                    height: width
                    radius: width/2
                    color: rectangleDelegate.border.color
                    opacity: parent.checked ? 1.0 : 0.0
                    Behavior on opacity { NumberAnimation { duration: 150; easing.type: Easing.OutQuint } }
                    scale: 0.5 + 0.5*opacity
                }

                ToolButton {
                    id: toolButtonChangeAccent
                    
                    anchors.fill: parent
                    anchors.margins: -6
                    z: 10

                    onClicked: {
                        currentSelectedIndex = index
                        applicationWindow.accentColor = Material.accent
                        saveConfiguration(Material.theme, Material.accentColor)
                    }
                } // ToolButton
            } // Rectangle
        } // Repeater
    } // Grid
}
