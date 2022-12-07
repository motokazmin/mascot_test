CREATE TABLE public.players (
	playername varchar NOT NULL,
	balance int4 NOT NULL,
	currency varchar(3) NOT NULL,
	CONSTRAINT players_un UNIQUE (playername)
);

CREATE TABLE public.transactions (
	transactionref varchar NOT NULL,
	playername varchar NOT NULL,
	gameid varchar NULL,
	sessionid varchar NULL,
	gameroundref varchar NULL,
	withdraw int4 NOT NULL,
	deposit int4 NOT NULL,
	currency varchar(3) NOT NULL,
	reason varchar NULL,
	bettype varchar NULL,
	wintype varchar NULL,
	id int4 NULL,
	CONSTRAINT transactions_un UNIQUE (transactionref)
);

INSERT INTO public.players (playername,balance,currency) VALUES
    ('player1',10000,'EUR'),
    ('player2',5000,'EUR');
