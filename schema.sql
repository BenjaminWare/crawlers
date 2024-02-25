-- +goose Down
-- DELETE ALL OF THE ABOVE TABLES
DROP TABLE IF EXISTS footnote_inst;

DROP TABLE IF EXISTS footnote;

DROP TABLE IF EXISTS non_derivative_transaction;

DROP TABLE IF EXISTS derivative_transaction;

DROP TABLE IF EXISTS form;

DROP TABLE IF EXISTS issuer;

DROP TABLE IF EXISTS stock_day;

DROP TABLE IF EXISTS ticker;

DROP TABLE IF EXISTS reporter;

-- +goose Up
CREATE TABLE issuer (
    cik varchar(10) primary key,
    created_at timestamp not null default current_timestamp,
    name varchar(200) not null,
    sic varchar(10),
    sic_description varchar(300),
    ein varchar(10),
    state_of_incorporation varchar(2),
    fiscal_year_end varchar(4),
    phone varchar(20),
    sector varchar(100),
    industry varchar(100)
);

CREATE TABLE ticker (
    id serial primary key,
    cik varchar(10) not null references issuer (cik),
    ticker varchar(10) not null,
    UNIQUE (ticker)
);

CREATE TABLE stock_day (
    id serial primary key,
    ticker varchar(10) not null references ticker (ticker),
    date varchar(10) not null,
    close decimal(19,4) not null
    -- Could include more info for each day, open high low ...
);

CREATE TABLE reporter (
    cik varchar(10) primary key,
    name varchar(200) not null
);

CREATE TABLE form (
    acc_num varchar(20) primary key,
    created_at timestamp not null default current_timestamp,
    period_of_report varchar(10) not null,
    rpt_is_director boolean not null,
    rpt_is_officer boolean not null,
    rpt_is_ten_percent_owner boolean not null,
    rpt_is_other boolean not null,
    rpt_officer_title varchar(100),
    rpt_other_text varchar(100),
    issuer_cik varchar(10) not null references issuer (cik),
    reporter_cik varchar(10) not null references reporter (cik),
    xml_url varchar(300) not null,
    pdf_url varchar(300) not null,
    net_shares decimal(19,4) not null,
    net_total decimal(19, 4) not null,
    transaction_codes varchar(20) not null
);

CREATE TABLE derivative_transaction (
    id serial primary key,
    acc_num varchar(20) not null references form (acc_num),
    security_title varchar(100),
    conversion_or_exercise_price decimal(19, 4),
    transaction_date varchar(10),
    transaction_form_type varchar(10),
    transaction_code varchar(10),
    equity_swap_involved boolean,
    transaction_shares decimal(19, 4),
    transaction_price_per_share decimal(19, 4),
    transaction_acquired_disposed_code varchar(10),
    exercise_date varchar(10),
    expiration_date varchar(10),
    underlying_security_title varchar(100),
    underlying_security_shares decimal(19, 4),
    post_transaction_amounts_shares decimal(19, 4),
    ownership_nature varchar(100),
    is_holding boolean not null
);

CREATE TABLE non_derivative_transaction(
    id serial primary key,
    acc_num varchar(20) not null references form (acc_num),
    security_title varchar(100),
    transaction_date varchar(10),
    transaction_form_type varchar(10),
    transaction_code varchar(10),
    equity_swap_involved boolean,
    transaction_shares decimal(19, 4),
    transaction_price_per_share decimal(19, 4),
    transaction_acquired_disposed_code varchar(10),
    post_transaction_amounts_shares decimal(19, 4),
    ownership_nature varchar(100),
    is_holding boolean not null
);

CREATE TABLE footnote(
    id serial primary key,
    acc_num varchar(20) not null references form (acc_num),
    text text not null
);

-- has surrugate key because instance may be on fields named the same, field_referenced isn't unique
CREATE TABLE footnote_inst(
    id serial primary key,
    acc_num varchar(20) not null references form (acc_num),
    footnote_id int not null references footnote(id),
    -- Joining footnote to footnote_inst requires both acc_num and footnote_id to footnote(id)
    -- At most one of these two attributes is not null, it indicates that a footnote references a field within a certain transaction, when joining make sure to include acc_num
    dt_id int references derivative_transaction(id),
    ndt_id int references non_derivative_transaction(id),
    field_referenced varchar(100) not null
);