const default_settings = {
    theme: "dark"
}
function main() {
    var settings = get_settings()
    apply_theme(settings.theme)
}
function apply_theme(theme) {
    var style = document.documentElement.style;
    switch (theme) {
        case "dark":
            style.setProperty("--base00", "#06070E")
            style.setProperty("--base01", "#0b0d14")
            style.setProperty("--text", "#f0f3ff")
            break;
        case "light":
            style.setProperty("--base00", "#f0f3ff")
            style.setProperty("--base01", "#e0e3ef")
            style.setProperty("--text", "#010103")
        default:
            break;
    }
}
function get_settings() {
    var settings = JSON.parse(window.localStorage.getItem("settingsData"))
    if (!settings) {
        set_settings(default_settings)
        set_settings = get_settings()
    }
    return settings
}

function set_settings(settings) {
    window.localStorage.setItem("settingsData", JSON.stringify(settings))
}
main()