# Codebase Backend

Backend of Codebase project.

## Developer Guide

### General Rules

* We use conventional commits to deal with git commits: <https://www.conventionalcommits.org>
  * Use `feat: commit message` to do git commit related to feature.
  * Use `refactor: commit message` to do git commit related to code refactorings.
  * Use `fix: commit message` to do git commit related to bugfix.
  * Use `test: commit message` to do git commit related to test files.
  * Use `docs: commit message` to do git commit related to documentations (including README.md files).
  * Use `style: commit message` to do git commit related to code style.

* Use git-chglog <https://github.com/git-chglog/git-chglog> to generate changelog (CHANGELOG.md) before merging to release branch.

### Branching Strategy

* Keep your branch strategy simple. Build your strategy from these three concepts:
  * Use feature branches for all new features and bug fixes.
  * Merge feature branches into the main branch using pull requests.
  * Keep a high quality, up-to-date main branch.

#### Use feature branches for your work

Develop your features and fix bugs in feature branches based off your main branch. These branches are also known as
topic branches. Feature branches isolate work in progress from the completed work in the main branch. Git branches are
inexpensive to create and maintain. Even small fixes and changes should have their own feature branch.

<!-- <p align="left"><img src="./featurebranching.png" width="360"></p> -->

#### Name your feature branches by convention

* Use a consistent naming convention for your feature branches to identify the work done in the branch. You can also
  include other information in the branch name, such as who created the branch.

* Some suggestions for naming your feature branches:
  * users/username/description
  * users/username/workitem
  * bugfix/description
  * feature/feature-name
  * feature/feature-area/feature-name
  * hotfix/description

#### Use release branches

* Create a release branch from the main branch when you get close to your release or other milestone, such as the end of
  a sprint. Give this branch a clear name associating it with the release, for example release/20.
* Create branches to fix bugs from the release branch and merge them back into the release branch in a pull request

<!-- <p align="left"><img src="./releasebranching_release.png" width="360"></p> -->

### Folder structure explanation

* `cmd/bin` folder is for storing the main.go file that will run the API server. this main.go file will call the `cmd/server` package to run the API server or with flag `seed` to seed the database with dummy data.
* `internal` folder is for storing the internal packages of the API server.
  * `adapter` folder is for storing the adapter struct which holds `driving adapters` and `driven adapters`.
    * **driving adapters** are the adapters that will be used in the API handler to interact with the service. e.g. Rest Server, CLI, Admin GUI.
    * **driven adapters** are the adapters that will be used in the service to interact with the repository.
  * `infrastucture` folder is for storing the infrastructure packages of the API server.
    * **config.go** is for storing the configuration the service needs.
    * **logging.go** is for storing the logger configuration.
  * `module` folder is for storing the modules of the API server that contains the entity, repository, service, and handler of a module.
    * `entity` is a package that contains the entity struct of a module.
    * `repository` is a package that contains the repository struct of a module.
    * `service` is a package that contains the service struct of a module.
    * `handler` is a package that contains the handler struct of a module.
  * `route` folder is for storing the route struct of the API server.
  * `logs` folder is for storing the log files of the API server.
  * `pkg` folder is for storing the common functions that will be used in the API server.

### How to create a new module

* there is a folder withing the `internal/module` folder named `z_template`.
* Copy the folder `z_template` and rename it to your module name.
* inside the folder, you will see folders named `entity`, `repository`, `service`, and `handler`.
  * `entity` folder is for storing the entity struct of your module. It is a struct that represents the data model (request and response commonly known as DTO and database model commonly known as DAO) of your module.
  * `repository` folder is for storing the repository struct of your module. It is **a** struct that contains the database operations of your module.
  * `service` folder is for storing the service struct of your module. It is a struct that contains the business logic of your module.
  * `handler` folder is for storing the handler struct of your module. It is a struct that contains the API handler of your module.
  * `ports` folder is for storing the ports of your module. It is a struct that contains the interfaces of your module.

### Write code step by step

 1. Read the requirements and create the structs in the `entity` folder (DTO and DAO).
 2. Create contract interfaces in the `ports` folder.
 3. Implement the interfaces in the `repository` folder.
 4. Implement the interfaces in the `service` folder.
 5. Create the handler in the `handler/rest` if it is a REST API. If it is a CLI, create the handler in the `handler/cli` folder and so on.
 6. Create injection in folder `handler/rest` if it is a REST API. If it is a CLI, create the injection in the `handler/cli` folder and so on.
