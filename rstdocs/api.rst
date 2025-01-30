
API
===


.. go:function:: func CreateClient(ctx context.Context, opts Options) (*Client, error)

    CreateClient returns a new client. The client connects lazily. Call
    Client.EnsureConnected() to force a connection.


.. go:function:: func CreateClientDSN(_ context.Context, dsn string, opts Options) (*Client, error)

    CreateClientDSN returns a new client. See also CreateClient.
    
    dsn is either an instance name
    `docs/clients/connection <https://www.edgedb.com/docs/clients/connection>`_
    or it specifies a single string in the following format:
    
    .. code-block:: go
    
        edgedb://user:password@host:port/database?option=value.
        
    The following options are recognized: host, port, user, database, password.


.. go:type:: type Client struct {\
        // contains filtered or unexported fields\
    }

    Client is a connection pool and is safe for concurrent use.


.. go:method:: func (p *Client) Close() error

    Close closes all connections in the pool.
    Calling close blocks until all acquired connections have been released,
    and returns an error if called more than once.


.. go:method:: func (p *Client) EnsureConnected(ctx context.Context) error

    EnsureConnected forces the client to connect if it hasn't already.


.. go:method:: func (p *Client) Execute(\
        ctx context.Context,\
        cmd string,\
        args ...interface{},\
    ) error

    Execute an EdgeQL command (or commands).


.. go:method:: func (p *Client) ExecuteSQL(\
        ctx context.Context,\
        cmd string,\
        args ...interface{},\
    ) error

    ExecuteSQL executes a SQL command (or commands).


.. go:method:: func (p *Client) Query(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    Query runs a query and returns the results.


.. go:method:: func (p *Client) QueryJSON(\
        ctx context.Context,\
        cmd string,\
        out *[]byte,\
        args ...interface{},\
    ) error

    QueryJSON runs a query and return the results as JSON.


.. go:method:: func (p *Client) QuerySQL(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    QuerySQL runs a SQL query and returns the results.


.. go:method:: func (p *Client) QuerySingle(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    QuerySingle runs a singleton-returning query and returns its element.
    If the query executes successfully but doesn't return a result
    a NoDataError is returned. If the out argument is an optional type the out
    argument will be set to missing instead of returning a NoDataError.


.. go:method:: func (p *Client) QuerySingleJSON(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    QuerySingleJSON runs a singleton-returning query.
    If the query executes successfully but doesn't have a result
    a NoDataError is returned.


.. go:method:: func (p *Client) Tx(ctx context.Context, action TxBlock) error

    Tx runs an action in a transaction retrying failed actions
    if they might succeed on a subsequent attempt.
    
    Retries are governed by retry rules.
    The default rule can be set with WithRetryRule().
    For more fine grained control a retry rule can be set
    for each defined RetryCondition using WithRetryCondition().
    When a transaction fails but is retryable
    the rule for the failure condition is used to determine if the transaction
    should be tried again based on RetryRule.Attempts and the amount of time
    to wait before retrying is determined by RetryRule.Backoff.
    If either field is unset (see RetryRule) then the default rule is used.
    If the object's default is unset the fall back is 3 attempts
    and exponential backoff.


.. go:method:: func (p Client) WithConfig(\
        cfg map[string]interface{},\
    ) *Client

    WithConfig sets configuration values for the returned client.


.. go:method:: func (p Client) WithGlobals(\
        globals map[string]interface{},\
    ) *Client

    WithGlobals sets values for global variables for the returned client.


.. go:method:: func (p Client) WithModuleAliases(\
        aliases ...ModuleAlias,\
    ) *Client

    WithModuleAliases sets module name aliases for the returned client.


.. go:method:: func (p Client) WithRetryOptions(\
        opts RetryOptions,\
    ) *Client

    WithRetryOptions returns a shallow copy of the client
    with the RetryOptions set to opts.


.. go:method:: func (p Client) WithTxOptions(opts TxOptions) *Client

    WithTxOptions returns a shallow copy of the client
    with the TxOptions set to opts.


.. go:method:: func (p Client) WithWarningHandler(\
        warningHandler WarningHandler,\
    ) *Client

    WithWarningHandler sets the warning handler for the returned client. If
    warningHandler is nil edgedb.LogWarnings is used.


.. go:method:: func (p Client) WithoutConfig(key ...string) *Client

    WithoutConfig unsets configuration values for the returned client.


.. go:method:: func (p Client) WithoutGlobals(globals ...string) *Client

    WithoutGlobals unsets values for global variables for the returned client.


.. go:method:: func (p Client) WithoutModuleAliases(\
        aliases ...string,\
    ) *Client

    WithoutModuleAliases unsets module name aliases for the returned client.


.. go:type:: type Error interface {\
        Error() string\
        Unwrap() error\
    \
        // HasTag returns true if the error is marked with the supplied tag.\
        HasTag(ErrorTag) bool\
    \
        // Category returns true if the error is in the provided category.\
        Category(ErrorCategory) bool\
    }

    Error is the error type returned from edgedb.


.. go:type:: type ErrorCategory string

    ErrorCategory values represent EdgeDB's error types.


.. go:type:: type ErrorTag string

    ErrorTag is the argument type to Error.HasTag().


.. go:type:: type Executor interface {\
        Execute(context.Context, string, ...any) error\
        Query(context.Context, string, any, ...any) error\
        QueryJSON(context.Context, string, *[]byte, ...any) error\
        QuerySingle(context.Context, string, any, ...any) error\
        QuerySingleJSON(context.Context, string, any, ...any) error\
    }

    Executor is a common interface between \*Client and \*Tx,
    that can run queries on an EdgeDB database.


.. go:type:: type IsolationLevel string

    IsolationLevel documentation can be found here
    `docs/reference/edgeql/tx_start#parameters <https://www.edgedb.com/docs/reference/edgeql/tx_start#parameters>`_


.. go:type:: type ModuleAlias struct {\
        Alias  string\
        Module string\
    }

    ModuleAlias is an alias name and module name pair.


.. go:type:: type Options struct {\
        // Host is an EdgeDB server host address, given as either an IP address or\
        // domain name. (Unix-domain socket paths are not supported)\
        //\
        // Host cannot be specified alongside the 'dsn' argument, or\
        // CredentialsFile option. Host will override all other credentials\
        // resolved from any environment variables, or project credentials with\
        // their defaults.\
        Host string\
    \
        // Port is a port number to connect to at the server host.\
        //\
        // Port cannot be specified alongside the 'dsn' argument, or\
        // CredentialsFile option. Port will override all other credentials\
        // resolved from any environment variables, or project credentials with\
        // their defaults.\
        Port int\
    \
        // Credentials is a JSON string containing connection credentials.\
        //\
        // Credentials cannot be specified alongside the 'dsn' argument, Host,\
        // Port, or CredentialsFile.  Credentials will override all other\
        // credentials not present in the credentials string with their defaults.\
        Credentials []byte\
    \
        // CredentialsFile is a path to a file containing connection credentials.\
        //\
        // CredentialsFile cannot be specified alongside the 'dsn' argument, Host,\
        // Port, or Credentials.  CredentialsFile will override all other\
        // credentials not present in the credentials file with their defaults.\
        CredentialsFile string\
    \
        // User is the name of the database role used for authentication.\
        //\
        // If not specified, the value is resolved from any compound\
        // argument/option, then from EDGEDB_USER, then any compound environment\
        // variable, then project credentials.\
        User string\
    \
        // Database is the name of the database to connect to.\
        //\
        // If not specified, the value is resolved from any compound\
        // argument/option, then from EDGEDB_DATABASE, then any compound\
        // environment variable, then project credentials.\
        //\
        // Deprecated: Database has been replaced by Branch\
        Database string\
    \
        // Branch is the name of the branch to use.\
        //\
        // If not specified, the value is resolved from any compound\
        // argument/option, then from EDGEDB_BRANCH, then any compound environment\
        // variable, then project credentials.\
        Branch string\
    \
        // Password to be used for authentication, if the server requires one.\
        //\
        // If not specified, the value is resolved from any compound\
        // argument/option, then from EDGEDB_PASSWORD, then any compound\
        // environment variable, then project credentials.\
        // Note that the use of the environment variable is discouraged\
        // as other users and applications may be able to read it\
        // without needing specific privileges.\
        Password types.OptionalStr\
    \
        // ConnectTimeout is used when establishing connections in the background.\
        ConnectTimeout time.Duration\
    \
        // WaitUntilAvailable determines how long to wait\
        // to reestablish a connection.\
        WaitUntilAvailable time.Duration\
    \
        // Concurrency determines the maximum number of connections.\
        // If Concurrency is zero, max(4, runtime.NumCPU()) will be used.\
        // Has no effect for single connections.\
        Concurrency uint\
    \
        // Parameters used to configure TLS connections to EdgeDB server.\
        TLSOptions TLSOptions\
    \
        // Read the TLS certificate from this file.\
        // DEPRECATED, use TLSOptions.CAFile instead.\
        TLSCAFile string\
    \
        // Specifies how strict TLS validation is.\
        // DEPRECATED, use TLSOptions.SecurityMode instead.\
        TLSSecurity string\
    \
        // ServerSettings is currently unused.\
        ServerSettings map[string][]byte\
    \
        // SecretKey is used to connect to cloud instances.\
        SecretKey string\
    \
        // WarningHandler is invoked when EdgeDB returns warnings. Defaults to\
        // edgedb.LogWarnings.\
        WarningHandler WarningHandler\
    }

    Options for connecting to an EdgeDB server


.. go:type:: type RetryBackoff func(n int) time.Duration

    RetryBackoff returns the duration to wait after the nth attempt
    before making the next attempt when retrying a transaction.


.. go:type:: type RetryCondition int

    RetryCondition represents scenarios that can cause a transaction
    run in Tx() methods to be retried.


.. go:function:: func NewRetryOptions() RetryOptions

    NewRetryOptions returns the default retry options.


.. go:type:: type RetryOptions struct {\
        // contains filtered or unexported fields\
    }

    RetryOptions configures how Tx() retries failed transactions.  Use
    NewRetryOptions to get a default RetryOptions value instead of creating one
    yourself.


.. go:method:: func (o RetryOptions) WithCondition(\
        condition RetryCondition,\
        rule RetryRule,\
    ) RetryOptions

    WithCondition sets the retry rule for the specified condition.


.. go:method:: func (o RetryOptions) WithDefault(rule RetryRule) RetryOptions

    WithDefault sets the rule for all conditions to rule.


.. go:function:: func NewRetryRule() RetryRule

    NewRetryRule returns the default RetryRule value.


.. go:type:: type RetryRule struct {\
        // contains filtered or unexported fields\
    }

    RetryRule determines how transactions should be retried when run in Tx()
    methods. See Client.Tx() for details.


.. go:method:: func (r RetryRule) WithAttempts(attempts int) RetryRule

    WithAttempts sets the rule's attempts. attempts must be greater than zero.


.. go:method:: func (r RetryRule) WithBackoff(fn RetryBackoff) RetryRule

    WithBackoff returns a copy of the RetryRule with backoff set to fn.


.. go:type:: type TLSOptions struct {\
        // PEM-encoded CA certificate\
        CA []byte\
        // Path to a PEM-encoded CA certificate file\
        CAFile string\
        // Determines how strict we are with TLS checks\
        SecurityMode TLSSecurityMode\
        // Used to verify the hostname on the returned certificates\
        ServerName string\
    }

    TLSOptions contains the parameters needed to configure TLS on EdgeDB
    server connections.


.. go:type:: type TLSSecurityMode string

    TLSSecurityMode specifies how strict TLS validation is.


.. go:type:: type Tx struct {\
        // contains filtered or unexported fields\
    }

    Tx is a transaction. Use Client.Tx() to get a transaction.


.. go:method:: func (t *Tx) Execute(\
        ctx context.Context,\
        cmd string,\
        args ...interface{},\
    ) error

    Execute an EdgeQL command (or commands).


.. go:method:: func (t *Tx) ExecuteSQL(\
        ctx context.Context,\
        cmd string,\
        args ...interface{},\
    ) error

    ExecuteSQL executes a SQL command (or commands).


.. go:method:: func (t *Tx) Query(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    Query runs a query and returns the results.


.. go:method:: func (t *Tx) QueryJSON(\
        ctx context.Context,\
        cmd string,\
        out *[]byte,\
        args ...interface{},\
    ) error

    QueryJSON runs a query and return the results as JSON.


.. go:method:: func (t *Tx) QuerySQL(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    QuerySQL runs a SQL query and returns the results.


.. go:method:: func (t *Tx) QuerySingle(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    QuerySingle runs a singleton-returning query and returns its element.
    If the query executes successfully but doesn't return a result
    a NoDataError is returned. If the out argument is an optional type the out
    argument will be set to missing instead of returning a NoDataError.


.. go:method:: func (t *Tx) QuerySingleJSON(\
        ctx context.Context,\
        cmd string,\
        out interface{},\
        args ...interface{},\
    ) error

    QuerySingleJSON runs a singleton-returning query.
    If the query executes successfully but doesn't have a result
    a NoDataError is returned.


.. go:type:: type TxBlock func(context.Context, *Tx) error

    TxBlock is work to be done in a transaction.


.. go:function:: func NewTxOptions() TxOptions

    NewTxOptions returns the default TxOptions value.


.. go:type:: type TxOptions struct {\
        // contains filtered or unexported fields\
    }

    TxOptions configures how transactions behave.


.. go:method:: func (o TxOptions) WithDeferrable(d bool) TxOptions

    WithDeferrable returns a shallow copy of the client
    with the transaction deferrable mode set to d.


.. go:method:: func (o TxOptions) WithIsolation(i IsolationLevel) TxOptions

    WithIsolation returns a copy of the TxOptions
    with the isolation level set to i.


.. go:method:: func (o TxOptions) WithReadOnly(r bool) TxOptions

    WithReadOnly returns a shallow copy of the client
    with the transaction read only access mode set to r.


.. go:type:: type WarningHandler = func([]error) error

    WarningHandler takes a slice of edgedb.Error that represent warnings and
    optionally returns an error. This can be used to log warnings, increment
    metrics, promote warnings to errors by returning them etc.


.. go:function:: func LogWarnings(errors []error) error

    LogWarnings is an edgedb.WarningHandler that logs warnings.


.. go:function:: func WarningsAsErrors(warnings []error) error

    WarningsAsErrors is an edgedb.WarningHandler that returns warnings as
    errors.