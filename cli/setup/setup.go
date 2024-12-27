// Package setup handles the interactive configuration of the Semaphore application
package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/semaphoreui/semaphore/util"
)

// interactiveSetupBlurb contains the welcome message and overview of the setup process
const interactiveSetupBlurb = `
Hello! You will now be guided through a setup to:

1. Set up configuration for a MySQL/MariaDB database. This step involves specifying the hostname, user, password, and database name for your MySQL instance.
2. Set up a path for your playbooks. A default path will be suggested, or you can provide your own location.
3. Run database migrations. This ensures that the database schema is up-to-date and ready for use.
4. Set up an initial semaphore user and password. This user will have administrative access to the application.
`

// InteractiveSetup guides the user through the configuration process
// It handles database selection, playbook path setup, alerts configuration, and LDAP setup
func InteractiveSetup(conf *util.ConfigType) {
	fmt.Print(interactiveSetupBlurb)

	selectDatabase(conf)
	setPlaybookPath(conf)
	configureAlerts(conf)
	configureLDAP(conf)
}

// selectDatabase prompts the user to choose between MySQL, BoltDB, or PostgreSQL
// and configures the selected database settings
func selectDatabase(conf *util.ConfigType) {
	dbPrompt := `What database to use:
   1 - MySQL
   2 - BoltDB
   3 - PostgreSQL
`
	var db int
	askValue(dbPrompt, "1", &db)

	switch db {
	case 1:
		conf.Dialect = util.DbDriverMySQL
		scanMySQL(conf)
	case 2:
		conf.Dialect = util.DbDriverBolt
		scanBoltDb(conf)
	case 3:
		conf.Dialect = util.DbDriverPostgres
		scanPostgres(conf)
	default:
		log.Warn("Invalid database selection. Defaulting to MySQL.")
		conf.Dialect = util.DbDriverMySQL
		scanMySQL(conf)
	}
}

// setPlaybookPath configures the path where playbooks will be stored
// Uses system temp directory as default base path
func setPlaybookPath(conf *util.ConfigType) {
	defaultPlaybookPath := filepath.Join(os.TempDir(), "semaphore")
	askValue("Playbook path", defaultPlaybookPath, &conf.TmpPath)
	conf.TmpPath = filepath.Clean(conf.TmpPath)
}

// configureAlerts sets up various alert integrations including email, Telegram,
// Slack, Rocket.Chat, and Microsoft Teams
func configureAlerts(conf *util.ConfigType) {
	alertConfigurations := []struct {
		name     string
		enabled  *bool
		settings func()
	}{
		{"email alerts", &conf.EmailAlert, func() {
			askValue("Mail server host", "localhost", &conf.EmailHost)
			askValue("Mail server port", "25", &conf.EmailPort)
			askValue("Mail sender address", "semaphore@localhost", &conf.EmailSender)
		}},
		{"telegram alerts", &conf.TelegramAlert, func() {
			askValue("Telegram bot token (from @BotFather)", "", &conf.TelegramToken)
			askValue("Telegram chat ID", "", &conf.TelegramChat)
		}},
		{"slack alerts", &conf.SlackAlert, func() {
			askValue("Slack Webhook URL", "", &conf.SlackUrl)
		}},
		{"Rocket.Chat alerts", &conf.RocketChatAlert, func() {
			askValue("Rocket.Chat Webhook URL", "", &conf.RocketChatUrl)
		}},
		{"Microsoft Teams alerts", &conf.MicrosoftTeamsAlert, func() {
			askValue("Microsoft Teams Webhook URL", "", &conf.MicrosoftTeamsUrl)
		}},
	}

	// Iterate through each alert type and configure if enabled
	for _, alert := range alertConfigurations {
		askConfirmation("Enable "+alert.name+"?", false, alert.enabled)
		if *alert.enabled {
			alert.settings()
		}
	}
}

// configureLDAP sets up LDAP authentication if enabled
// Configures server details, bind credentials, and field mappings
func configureLDAP(conf *util.ConfigType) {
	askConfirmation("Enable LDAP authentication?", false, &conf.LdapEnable)
	if conf.LdapEnable {
		conf.LdapMappings = &util.LdapMappings{}
		askValue("LDAP server host", "localhost:389", &conf.LdapServer)
		askConfirmation("Enable LDAP TLS connection", false, &conf.LdapNeedTLS)
		askValue("LDAP DN for bind", "cn=user,ou=users,dc=example", &conf.LdapBindDN)
		askValue("Password for LDAP bind user", "", &conf.LdapBindPassword)
		askValue("LDAP DN for user search", "ou=users,dc=example", &conf.LdapSearchDN)
		askValue("LDAP search filter", `(uid=%s)`, &conf.LdapSearchFilter)
		askValue("LDAP mapping for DN field", "dn", &conf.LdapMappings.DN)
		askValue("LDAP mapping for username field", "uid", &conf.LdapMappings.UID)
		askValue("LDAP mapping for full name field", "cn", &conf.LdapMappings.CN)
		askValue("LDAP mapping for email field", "mail", &conf.LdapMappings.Mail)
	}
}

// scanBoltDb configures BoltDB settings
// Uses the current working directory for the default database file location
func scanBoltDb(conf *util.ConfigType) {
	workingDirectory := getWorkingDirectory()
	defaultBoltDBPath := filepath.Join(workingDirectory, "database.boltdb")
	conf.BoltDb = &util.DbConfig{}
	askValue("db filename", defaultBoltDBPath, &conf.BoltDb.Hostname)
}

// scanMySQL configures MySQL connection settings
// Sets up hostname, user, password, and database name
func scanMySQL(conf *util.ConfigType) {
	conf.MySQL = &util.DbConfig{}
	askValue("db Hostname", "127.0.0.1:3306", &conf.MySQL.Hostname)
	askValue("db User", "root", &conf.MySQL.Username)
	askValue("db Password", "", &conf.MySQL.Password)
	askValue("db Name", "semaphore", &conf.MySQL.DbName)
}

// scanPostgres configures PostgreSQL connection settings
// Sets up hostname, user, password, database name, and SSL mode
func scanPostgres(conf *util.ConfigType) {
	conf.Postgres = &util.DbConfig{}
	askValue("db Hostname", "127.0.0.1:5432", &conf.Postgres.Hostname)
	askValue("db User", "root", &conf.Postgres.Username)
	askValue("db Password", "", &conf.Postgres.Password)
	askValue("db Name", "semaphore", &conf.Postgres.DbName)
	if conf.Postgres.Options == nil {
		conf.Postgres.Options = make(map[string]string)
	}
	conf.Postgres.Options["sslmode"] = "disable"
}

// getWorkingDirectory returns the current working directory
// Falls back to system temp directory if working directory cannot be determined
func getWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return os.TempDir()
	}
	return workingDirectory
}

// askValue prompts the user for input with an optional default value
// The provided item pointer is updated with the user's input
func askValue(prompt, defaultValue string, item interface{}) {
	fmt.Print(prompt)
	if defaultValue != "" {
		fmt.Printf(" (default %s)", defaultValue)
	}
	fmt.Print(": ")
	_, _ = fmt.Sscanln(defaultValue, item)
	scanErrorChecker(fmt.Scanln(item))
	fmt.Println("")
}

// askConfirmation prompts the user for a yes/no response
// Updates the provided boolean pointer based on the user's input
func askConfirmation(prompt string, defaultValue bool, item *bool) {
	defString := "yes"
	if !defaultValue {
		defString = "no"
	}
	fmt.Printf("%s (yes/no) (default %s): ", prompt, defString)
	var answer string
	scanErrorChecker(fmt.Scanln(&answer))
	switch strings.ToLower(answer) {
	case "y", "yes":
		*item = true
	case "n", "no":
		*item = false
	default:
		*item = defaultValue
	}
	fmt.Println("")
}

// scanErrorChecker handles errors that occur during user input scanning
// Ignores expected newline errors but logs other scanning errors
func scanErrorChecker(err error) {
	if err != nil && err.Error() != "unexpected newline" {
		log.Warn("An input error occurred: " + err.Error())
	}
}