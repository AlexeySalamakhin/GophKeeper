// Package client содержит CLI клиент GophKeeper.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Client представляет CLI клиент.
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	configPath string
}

// New создает новый экземпляр клиента.
func New() *Client {
	return &Client{
		baseURL: viper.GetString("server.url"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		configPath: viper.GetString("config.path"),
	}
}

// Execute запускает CLI клиент.
func (c *Client) Execute() error {
	c.loadToken()

	rootCmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "GophKeeper - менеджер паролей",
		Long:  "GophKeeper - безопасный менеджер паролей с синхронизацией",
	}

	// Команды аутентификации
	rootCmd.AddCommand(c.createAuthCommands())

	// Команды для работы с данными
	rootCmd.AddCommand(c.createDataCommands())

	// Команда версии
	rootCmd.AddCommand(c.createVersionCommand())

	return rootCmd.Execute()
}

// createAuthCommands создает команды аутентификации.
func (c *Client) createAuthCommands() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Команды аутентификации",
	}

	// Команда регистрации
	registerCmd := &cobra.Command{
		Use:   "register [username] [email] [password]",
		Short: "Регистрация нового пользователя",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			c.register(args[0], args[1], args[2])
		},
	}

	// Команда входа
	loginCmd := &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Вход в систему",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c.login(args[0], args[1])
		},
	}

	// Команда выхода
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Выход из системы",
		Run: func(cmd *cobra.Command, args []string) {
			c.logout()
		},
	}

	authCmd.AddCommand(registerCmd, loginCmd, logoutCmd)
	return authCmd
}

// createDataCommands создает команды для работы с данными.
func (c *Client) createDataCommands() *cobra.Command {
	dataCmd := &cobra.Command{
		Use:   "data",
		Short: "Команды для работы с данными",
	}

	// Команда списка данных
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Список всех данных",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			c.listData()
		},
	}

	// Команда добавления данных
	addCmd := &cobra.Command{
		Use:   "add [name] [login] [password]",
		Short: "Добавить новые данные",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			metadata, _ := cmd.Flags().GetString("metadata")
			c.addData(args[0], args[1], args[2], metadata)
		},
	}
	addCmd.Flags().String("metadata", "", "Метаданные в формате JSON")

	// Команда получения данных
	getCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Получить данные по ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.getData(args[0])
		},
	}

	dataCmd.AddCommand(listCmd, addCmd, getCmd)
	return dataCmd
}

// createVersionCommand создает команду версии.
func (c *Client) createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Показать версию и дату сборки",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("GophKeeper v1.0.0")
			fmt.Println("Дата сборки:", time.Now().Format("2006-01-02 15:04:05"))
		},
	}
}

// register выполняет регистрацию пользователя.
func (c *Client) register(username, email, password string) {
	req := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}

	resp, err := c.makeRequest("POST", "/api/v1/register", req)
	if err != nil {
		fmt.Printf("Ошибка регистрации: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var authResp struct {
			Token string `json:"token"`
			User  struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"user"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&authResp); err == nil {
			c.token = authResp.Token
			if err := c.saveToken(); err != nil {
				fmt.Printf("Предупреждение: не удалось сохранить токен: %v\n", err)
			}
			fmt.Printf("Успешная регистрация! Добро пожаловать, %s!\n", authResp.User.Username)
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Ошибка регистрации: %s\n", string(body))
	}
}

// login выполняет вход пользователя.
func (c *Client) login(username, password string) {
	req := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := c.makeRequest("POST", "/api/v1/login", req)
	if err != nil {
		fmt.Printf("Ошибка входа: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var authResp struct {
			Token string `json:"token"`
			User  struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"user"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&authResp); err == nil {
			c.token = authResp.Token
			if err := c.saveToken(); err != nil {
				fmt.Printf("Предупреждение: не удалось сохранить токен: %v\n", err)
			}
			fmt.Printf("Успешный вход! Добро пожаловать, %s!\n", authResp.User.Username)
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Ошибка входа: %s\n", string(body))
	}
}

// logout выполняет выход пользователя.
func (c *Client) logout() {
	c.token = ""
	if err := c.saveToken(); err != nil {
		fmt.Printf("Предупреждение: не удалось сохранить токен: %v\n", err)
	}
	fmt.Println("Выход выполнен успешно")
}

// listData выводит список данных.
func (c *Client) listData() {
	if c.token == "" {
		fmt.Println("Необходимо войти в систему")
		return
	}

	resp, err := c.makeRequest("GET", "/api/v1/data", nil)
	if err != nil {
		fmt.Printf("Ошибка получения данных: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			if len(data) == 0 {
				fmt.Println("Данные не найдены")
				return
			}

			fmt.Printf("Найдено %d записей:\n", len(data))
			for _, item := range data {
				fmt.Printf("- ID: %v, Название: %v, Логин: %v\n",
					item["id"], item["name"], item["login"])
			}
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Ошибка получения данных: %s\n", string(body))
	}
}

// addData добавляет новые данные.
func (c *Client) addData(name, login, password, metadata string) {
	if c.token == "" {
		fmt.Println("Необходимо войти в систему")
		return
	}

	// Сборка запроса
	req := map[string]interface{}{
		"name":     name,
		"login":    login,
		"password": password,
	}
	if metadata != "" {
		var metaObj map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &metaObj); err == nil {
			req["metadata"] = metaObj
		} else {
			fmt.Println("Предупреждение: метаданные должны быть JSON, игнорируются")
		}
	}

	resp, err := c.makeRequest("POST", "/api/v1/data", req)
	if err != nil {
		fmt.Printf("Ошибка добавления данных: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("Данные успешно добавлены")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Ошибка добавления данных: %s\n", string(body))
	}
}

// getData получает данные по ID.
func (c *Client) getData(id string) {
	if c.token == "" {
		fmt.Println("Необходимо войти в систему")
		return
	}

	resp, err := c.makeRequest("GET", "/api/v1/data/"+id, nil)
	if err != nil {
		fmt.Printf("Ошибка получения данных: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			fmt.Printf("Данные: %+v\n", data)
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Ошибка получения данных: %s\n", string(body))
	}
}

// makeRequest выполняет HTTP запрос.
func (c *Client) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.httpClient.Do(req)
}

// saveToken сохраняет токен в файл.
func (c *Client) saveToken() error {
	if err := os.MkdirAll(c.configPath, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию конфигурации: %w", err)
	}

	tokenFile := filepath.Join(c.configPath, "token")
	if err := os.WriteFile(tokenFile, []byte(c.token), 0600); err != nil {
		return fmt.Errorf("не удалось сохранить токен: %w", err)
	}

	return nil
}

// loadToken загружает токен из файла.
func (c *Client) loadToken() {
	tokenFile := filepath.Join(c.configPath, "token")
	if data, err := os.ReadFile(tokenFile); err == nil {
		c.token = string(data)
	}
}
