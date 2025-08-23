package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Danyarbrg/snippetbox/pkg/models/mysql"
	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql"
)

// определяем структуру содержащую зависимости для веб приложения.
// мы включаем поля для кастомных логгеров
// также добавили поле SnippetModel чтобы его можно было использовать в наших обработчиках
type application struct {
	errorLog 		*log.Logger
	infoLog			*log.Logger
	session			*sessions.Session
	snippets 		*mysql.SnippetModel
	templateCache	map[string]*template.Template
	users			*mysql.UserModel
}

func main() {
	// определяем новый терминальный флаг с именем 'addr', дефолтное значение которого 
	// ":4000" и немного короткого текста для понимания что флаг контролирует
	// флаг будет расположен в addr заначение во время запуска команды
	// bash: go run ./cmd/web -addr=":9999"
	// также можно использовать: go run ./cmd/web -help чтобы увидеть какие переменные задейстованны
	addr := flag.String("addr", ":4000", "HTTP network addres")
	// определяем новый флаг для MySQL DSN 
	dsn := flag.String("dsn", "web:1234@/snippetboxdb?parseTime=true", "MySQL data source name")
	// define new command-line flag for the session secret (random key)
	// it should be 32 bytes long
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key") 
	// мы вызываем flag.Parse() читает значение и передает его в терминал.
	// нужно вызывать перед использованием addr значения, иначе ему будет
	// присвоено стандартное значение ":4000"
	flag.Parse()

	// log.New() создает новый логгер для информационного сообщеня. Он принимает три парамтера:
	// первый - пункт назначения для записи журналов, в нашем случае вывод в терминал
	// второй - префикс для выводимого сообщения
	// третий - и флаги с дополнительной ифномрацией. несколько флагов связаны оператором OR а именно |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// log.New() создает новый логгер для ошибочного сообщения. 
	// но мы используем stderr как пункт назначения и добавили флаг 
	// Lshortfile для указания конркетного файла и номера строки
	// Llongfile указывать полный путь до файла
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// use sessions.New() func to inizialize a new session manager
	// passing in the secret key as the parametr. 
	// then session always expires after 12 hours
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// инициализиурем новый экземпляр приложения содержащий зависимости
	app := &application{
		errorLog: 		errorLog,
		infoLog: 		infoLog,
		session: 		session,
		snippets:		&mysql.SnippetModel{DB: db},
		templateCache:	templateCache,
		users:			&mysql.UserModel{DB: db},
	}

	// inicialize a tls.Config struct to hold the non-defaults TLS settings we want the server to use
	tlsConfig := &tls.Config{
		PreferServerCipherSuites:	true,
		CurvePreferences:			[]tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// инициализировали новой http.Server структуры. В поля Addr и Handler 
	// мы поместили данные которые у нас уже до этого использовались. А в поле ErrorLog
	// мы указали кастомное логирование ошибок
	srv := &http.Server{
		Addr:		*addr,
		ErrorLog: 	errorLog,
		Handler: 	app.routes(),
		TLSConfig: 	tlsConfig,
		// add idle, read and write timeouts to the server
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// значение возвращаемое из flag.String() функции это указатель на флаг-значение,
	// не на само значение. По этому нам надо раскрыть значение перед его использованием
	// infoLog - пишет новые созданные логи вместо обычных
	// go run ./cmd/web >>./logs/info.log 2>>./logs/error.log -- записывать логи в файлы
	// где > -- перезаписывать. >> -- дописывать в конец файла
	infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux) -- изначальный вариант
	// вариант с использованием метода ListenAndServe для нашей новой структуры
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}