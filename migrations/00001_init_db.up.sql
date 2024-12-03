BEGIN;

CREATE TABLE criterios_admissao (
  id integer NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  nome varchar NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE criterios_admissao_servico (
  servico_id integer NOT NULL,
  criterio_admissao_id integer NOT NULL
);

CREATE TABLE forma_encaminhamento (
  id integer NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  nome varchar NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE forma_encaminhamento_servico (
  forma_encaminhamento_id integer NOT NULL,
  servico_id integer NOT NULL
);

CREATE TABLE regionais (
  id integer NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  nome varchar NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE servico (
  id integer NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  tipo_servico_id INTEGER NOT NULL,
  nome varchar NOT NULL,
  endereco varchar NOT NULL,
  contato varchar,
  regional varchar,
  site varchar,
  regional_id integer NOT NULL,
  observacoes varchar,
  PRIMARY KEY (id)
);

CREATE TABLE tipo_atendimento (
  id integer NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  nome varchar NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE tipo_atendimento_servico (
  tipo_atendimento_id integer NOT NULL,
  servico_id integer NOT NULL
);

CREATE TABLE tipo_servico (
  id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  nome VARCHAR NOT NULL,
  PRIMARY KEY (id)
);

ALTER TABLE
  servico
ADD
  CONSTRAINT FK_tipo_servico_TO_servico FOREIGN KEY (tipo_servico_id) REFERENCES tipo_servico (id);

ALTER TABLE
  servico
ADD
  CONSTRAINT FK_regionais_TO_servico FOREIGN KEY (regional_id) REFERENCES regionais (id);

ALTER TABLE
  tipo_atendimento_servico
ADD
  CONSTRAINT FK_tipo_atendimento_TO_tipo_atendimento_servico FOREIGN KEY (tipo_atendimento_id) REFERENCES tipo_atendimento (id);

ALTER TABLE
  tipo_atendimento_servico
ADD
  CONSTRAINT FK_servico_TO_tipo_atendimento_servico FOREIGN KEY (servico_id) REFERENCES servico (id);

ALTER TABLE
  criterios_admissao_servico
ADD
  CONSTRAINT FK_servico_TO_criterios_admissao_servico FOREIGN KEY (servico_id) REFERENCES servico (id);

ALTER TABLE
  criterios_admissao_servico
ADD
  CONSTRAINT FK_criterios_admissao_TO_criterios_admissao_servico FOREIGN KEY (criterio_admissao_id) REFERENCES criterios_admissao (id);

ALTER TABLE
  forma_encaminhamento_servico
ADD
  CONSTRAINT FK_forma_encaminhamento_TO_forma_encaminhamento_servico FOREIGN KEY (forma_encaminhamento_id) REFERENCES forma_encaminhamento (id);

ALTER TABLE
  forma_encaminhamento_servico
ADD
  CONSTRAINT FK_servico_TO_forma_encaminhamento_servico FOREIGN KEY (servico_id) REFERENCES servico (id);

COMMIT;