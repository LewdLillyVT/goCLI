# For Developers - Adding Plugins to goCLI

Welcome to the goCLI plugin development guide! This document explains how to enhance your JavaScript plugins by adding Node.js modules and how to document your plugins with code comments.

## Adding Node.js Modules

To use Node.js modules in your JavaScript plugins, you can include them at the top of your `.js` file. Make sure to install the necessary modules before using them in your code. Here’s how you can do it:

1. **Install Node Modules**: Use npm (Node Package Manager) to install any required modules. For example, if you want to use the `readline-sync` module, you would run:

   ```bash
   npm install readline-sync
   ```

2. **Include Modules in Your Plugin**: At the top of your JavaScript file, you can specify the required modules in comments and then require them. Here’s how you can do it:

   ```javascript
   // require: readline-sync

   // Your actual plugin code starts here
   ```

Here’s an example of how the top of your JavaScript plugin might look:

```javascript
// require: readline-sync

// Your actual plugin code starts here
function factorial(n) {
    if (n < 0) {
        return "Invalid input. Please enter a non-negative integer.";
    }
    var result = 1; 
    for (var i = 1; i <= n; i++) {
        result *= i;
    }
    return result;
}

function run(getInputFunction) {
    // Using the getInputFunction to get user input
    var input = getInputFunction("Enter a number to calculate its factorial:"); 
    var number = parseInt(input); 

    if (isNaN(number)) {
        console.log("Invalid input. Please enter a valid number.");
        return "Invalid input. Please enter a valid number.";
    }

    var result = factorial(number); 
    var outputMessage = "The factorial of " + number + " is " + result + "."; 
    console.log(outputMessage);
    return outputMessage;
}

// factorial.js by LewdLillyVT
// https://github.com/LewdLillyVT/PluginLib
// An example plugin to test the new node package manager integration feature
```

## Documenting Plugin Information

To ensure that users of your plugin have the necessary information, please add three lines of code comments at the end of each plugin file. This information will be read by the goCLI application to provide users with details about the plugin.

### Format of the Plugin Information Comments

The last three lines should ideally include:

1. A brief description of what the plugin does.
2. Any specific usage instructions or requirements.
3. The author's name or a link to the plugin repository.

### Example

At the end of your JavaScript plugin, you might add:

```javascript
// Author: LewdLillyVT
// This plugin calculates the factorial of a number using user input.
// Ensure you have the necessary Node.js modules installed.
```

Thank you for contributing to goCLI!
