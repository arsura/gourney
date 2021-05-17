CREATE TABLE currencies (
  id bigserial PRIMARY KEY,
  name character varying NOT NULL,
  amount double precision,
  total double precision,
  rise_rate double precision,
  rise_factor double precision,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE exchange_rates (
  id bigserial PRIMARY KEY,
  source_id bigint NULL REFERENCES currencies(id),
  target_id bigint NULL REFERENCES currencies(id),
  price double precision,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE histories (
  id bigserial PRIMARY KEY,
  client_id bigint NOT NULL,
  source_id bigint NULL REFERENCES currencies(id),
  target_id bigint NULL REFERENCES currencies(id),
  sorce_amount double precision,
  target_amount double precision,
  price double precision,
  created_at timestamptz NOT NULL DEFAULT (now())
);