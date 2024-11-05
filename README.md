# goCLI

**goCLI** is a modular CLI tool written in **Go**. It is designed to support a flexible plugin system, allowing users to extend functionality with **Python**, **JavaScript**, and **PowerShell** plugins. All plugins are community-generated, making goCLI highly adaptable to different use cases.

---

## 🔌 Plugins

Plugins in goCLI are user-generated code snippets that interact seamlessly with the goCLI environment. They allow users to add custom features and functionality, making goCLI highly customizable. goCLI provides a simple way to install plugins via its native **[Plugin Library](https://github.com/YOUR_PLUGIN_LIB_REPO_LINK)**, where community members can submit their own plugins for review and potential inclusion.

The plugin review process ensures a high standard of quality and security. However, users and developers can also manually add plugins locally, thanks to goCLI's structured folder system.

---

## 📂 Folder Structure

Upon first run, goCLI will create a new directory at `/LOCALAPPDATA/goCLI` with the following subfolders:

- **logs**: Stores error logs for debugging purposes.
- **plugins**: Stores downloaded plugins. Users can also add their own plugins directly to this folder.
- **dependencies**: Exclusively used for JavaScript plugins that require external modules. This is necessary because our JavaScript interpreter, Otto, does not natively support NPM modules.

---

## ⚠️ Important Note

On a fresh Windows installation, goCLI will only support PowerShell plugins and JavaScript plugins that don't require additional modules. To unlock full functionality:

1. **Install Python 3** to enable Python plugins.
2. **Install Node.js and NPM** for JavaScript plugins that require external modules.

Ensure that both Python and NPM are added to your system PATH.

---

Thank you for using goCLI! We hope this tool and the contributions from the community make goCLI an invaluable addition to your toolkit.
