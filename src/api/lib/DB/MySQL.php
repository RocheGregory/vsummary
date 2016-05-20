<?php namespace Ozdemir\Datatables\DB;

use PDO;
use PDOException;

class MySQL implements DatabaseInterface {

    private $pdo;
    private $config;
    private $escape = [];

    function __construct($config)
    {
        $this->config = $config;
    }

    public function connect()
    {
        $host = $this->config['host'];
        $port = $this->config['port'];
        $user = $this->config['username'];
        $pass = $this->config['password'];
        $database = $this->config['database'];
        $charset = 'utf8';

        try {
            $this->pdo = new PDO("mysql:host=$host;dbname=$database;port=$port;charset=$charset", "$user", "$pass");
        } catch ( PDOException $e ){
            print $e->getMessage();
        }

        return $this;
    }

    public function query($query)
    {
        $sql = $this->pdo->prepare($query);
        $rows=$sql->execute($this->escape);

        return $sql->fetchAll(PDO::FETCH_ASSOC);
    }

    public function count($query)
    {
        $sql = $this->pdo->prepare($query);
        $rows=$sql->execute($this->escape);

        return $sql->rowCount();
    }

    public function escape($string)
    {
        $this->escape[':escape' . (count($this->escape) + 1) ] = '%' . $string . '%';

        return ":escape" . (count($this->escape));
    }

}