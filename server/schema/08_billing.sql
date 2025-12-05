\c lesnotes
CREATE SCHEMA IF NOT EXISTS billing;

CREATE TYPE billing.currency AS ENUM ('rub', 'usd', 'eur', 'ton', 'btc', 'xtr');
CREATE TYPE billing.invoice_status AS ENUM ('paid', 'unpaid', 'cancelled');
CREATE TYPE billing.payment_status AS ENUM ('refunded', 'cancelled', 'processed', 'pending');

CREATE TABLE IF NOT EXISTS billing.payments
(
	id             bigint            UNIQUE NOT NULL,
	invoice_id     VARCHAR(256)      NOT NULL,
	user_id        bigint            NOT NULL,
	status         billing.payment_status        NOT NULL DEFAULT 'pending',
	currency       billing.currency  NOT NULL DEFAULT 'rub',
	total          bigint            NOT NULL,
	created_at     timestamptz       NOT NULL DEFAULT NOW(),
	updated_at     timestamptz       NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_billing_payments_trgr BEFORE UPDATE ON billing.payments FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_billing_payments_trgr BEFORE UPDATE ON billing.payments FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS billing.telegram
(
	id             bigint            NOT NULL,
	user_id        bigint            NOT NULL,
	invoice_id     VARCHAR(256)      NOT NULL,
	status         billing.payment_status        NOT NULL DEFAULT 'pending',
	telegram_payment_charge_id     text          NOT NULL,
	provider_payment_charge_id     text          NOT NULL,
	currency       billing.currency  NOT NULL,
	total_amount   integer           NOT NULL,
	created_at     timestamptz       NOT NULL DEFAULT NOW(),
	updated_at     timestamptz       NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_billing_telegram_trgr BEFORE UPDATE ON billing.telegram FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_billing_telegram_trgr BEFORE UPDATE ON billing.telegram FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS billing.invoices
(
	id             VARCHAR(256)            UNIQUE NOT NULL,
	user_id        bigint                  NOT NULL,
	status         billing.invoice_status  NOT NULL DEFAULT 'unpaid',
	currency       billing.currency        NOT NULL DEFAULT 'rub',
	total          bigint                  NOT NULL,
	created_at     timestamptz             NOT NULL DEFAULT NOW(),
	updated_at     timestamptz             NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_billing_invoices_trgr BEFORE UPDATE ON billing.invoices FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_billing_invoices_trgr BEFORE UPDATE ON billing.invoices FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA billing TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA billing TO lesnotes_admin;

GRANT USAGE ON TYPE billing.currency TO lesnotes_admin;
GRANT USAGE ON TYPE billing.invoice_status TO lesnotes_admin;
GRANT USAGE ON TYPE billing.payment_status TO lesnotes_admin;
