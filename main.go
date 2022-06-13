package main

import (
	"fmt"

	client "github.com/aptible/cloud-api-clients"
)

type LoginInput struct {
	email string
	pass  string
}

type CreateRDSParams struct {
	Name          string // e.g. name of db
	Engine        string // e.g. postgres
	EngineVersion string // e.g. 14.2
}

/*
The goal of this interface is to be an abstraction layer above the cloud-api.
Whenever we want to interface with the API, we should use this infterface.
*/
type CloudClient interface {
	ListEnvironments() error
	CreateEnvironment(orgID string, params *client.EnvironmentInput) error
	RemoveEnvironment(orgID string, envID string) error

	ListAssetTypesForEnvironment(envID string) error
	CreateAsset(orgID string, envID string, params *client.AssetInput) error
	UpdateAsset(orgID string, envID string, assetID string, assetInput *client.AssetInput) error
	DeleteAsset(orgID string, envID string, assetID string) error
}

/*
The goal of this interface is to represent the commands that we plan to implement.
Whatever CLI framework we use will primarily interact with this interface.
*/
type CLI interface {
	Login(params *LoginInput) error
	Version() string

	ListOrgs() error
	ListEnvironments(orgID string) error
	ListApps(orgID string) error
	ListDatastores(orgID string) error
	ListsLogs(orgID string) error
	ListLogsForAsset(orgID string, envID string, assetID string) error

	SSH(orgID string, envID string, assetID string) error
	Status(orgID string, envID string, assetID string) error
	Info(orgID string, envID string, assetID string) error
	Open(orgID string, envID string, assetID string) error

	CreateRDS(input *CreateRDSParams) error

	CreateBackup(orgID string, envID string, assetID string) error
	DeleteDatastore(orgID string, envID string, assetID string) error
	DeleteBackup(orgID string, envID string, backupID string) error
	RestoreDatastore(backupID string, orgID string, envID string, assetID string) error
}

func main() {
	fmt.Println("Hello world!")
}

/*
HELP
----

aptible is a command line interface to the Aptible.com platform.

It allows users to manage authentication, application launch,
deployment, logging, and more with just the one command.

* Deploy an app with the app deploy command
* Provision a datastore with the datastore deploy command
* View a deployed web application with the open command
* View detailed information about an app or datastore with the info command

To read more, use the docs command to view Aptible's help on the web.

Usage:
	aptible [command]

Available commands:
	app			Manage your Aptible applications
	auth 		Manage authentication
	backup		Manage datastore backups
	datastore	Manage your Aptible datastores
	doc			View Aptible documentation
	env 		Manage environments
	help 		Help about any command
	info 		Show detailed info about an app or datastore
	list		List your Aptible resources
	log			Tail logs for an app or datastore
	open 		Open browser to current deployed application
	org 		Manage organizations
	ssh			SSH into an app
*/

/*
DATASTORE
---------

The datastore subcommand helps manage your Aptible datastores.

Usage:
	aptible datastore [command]

Aliases:
	ds, db

Available commands:
	list	list all datastores
	deploy	creates a new datastore
	destroy	destroys a datastore
	tunnel	create a tunnel from a datastore to your local machine

Flags:
	-h	help for datastore command

Use "aptible datastore [command] --help" for more information about a command.
*/

/*
DATASTORE DEPLOY
----------------

The datastore deploy command will provision a new datastore.

Usage:
	aptible datastore deploy [ENVIRONMENT]/[NAME]

Flags:
	-h				help for deploy
	-e, --engine	the datastore engine, e.g. postgres, mysql, etc.
	-v, --version	the engine version, e.g. 14.2
*/

/*
DATASTORE DESTROY
-----------------

The datastore destroy command will permentantly remove the datastore.

Usage:
	aptible datastore destroy [ENVIRONMENT]/[NAME]

Flags:
	-h	help for destroy
*/

/*
BACKUP
------

The backup subcommand helps manage your Aptible backups.

Usage:
	aptible backup [command]

Aliases:
	b

Available commands:
	create	creates a backup from a datastore
	destroy	deletes a backup
	restore	uses backup to deploy a datastore

Flags:
	-h	help for backup command

Use "aptible backup [command] --help" for more information about a command.
*/

/*
BACKUP CREATE
-------------

The backup create command will take a snapshot of the datastore provided.

Usage:
	aptible backup create [ENVIRONMENT]/[NAME]

Flags:
	-h			help for backup create command
	-l, --label	the label for the backup
*/

/*
BACKUP DESTROY
--------------

The backup destroy command will permentantly delete a backup.

Usage:
	aptible backup destroy (label|uuid)

Flags:
	-h			help for backup destroy command
*/

/*
BACKUP RESTORE
--------------

The backup restore command will create a datastore from a backup.

Usage:
	aptible backup restore (label|uuid) [ENVIRONMENT]/[NAME]

Flags:
	-h			help for backup restore command
*/
