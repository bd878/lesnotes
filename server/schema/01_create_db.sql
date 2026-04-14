CREATE DATABASE lesnotes;

GRANT CONNECT ON DATABASE lesnotes TO lesnotes_admin;
GRANT ALL PRIVILEGES ON DATABASE lesnotes TO lesnotes_admin;
ALTER ROLE lesnotes_admin SET search_path TO lesnotes, "$user", public;

\c lesnotes
