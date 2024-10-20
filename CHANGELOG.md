# Changelog
## Alpha-Release 0.3.0 `20.10.2024`
* **Important**: This is a breaking release. Due to the changed format of the parts from base64 to []bytes there will be compatibility issues with older version.
* Completely removing base64 as transport format and working with []byte directly between zip and AES GCM. ( https://github.com/voodooEntity/go-tachicrypt/issues/3 )
* Adjusting cli output to remove any traces of base64 information ( https://github.com/voodooEntity/go-tachicrypt/issues/3 )
* Adjusting cli output configuration part to name the given args "path" instead of "directory" since it can either be a file or a directory
* Some code/comment cleanup
* Restructuring and adjusting README.md to changes 

## Alpha-Release 0.2.0 `19.10.2024`
* Small update of README.md
* Adding prettywriter package for cleaner cli output
* Massive overhaul of cli output
* Some code/comment cleanup
* Some error handling improvements

## Alpha-Release 0.1.0 `18.10.2024`
* Adding CHANGELOG.md
* Adjusting the password input method to golang.org/x/term for non visible password input
* Adjust .gitignore
* Some output cleanup

## Alpha-Release 0.0.1 `13.10.2024`
* First Alpha release of tachicrypt