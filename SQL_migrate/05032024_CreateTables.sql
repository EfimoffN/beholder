
CREATE TABLE if not exists prj_client(
    clientid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    clientname CHARACTER VARYING(300) NOT NULL UNIQUE,
    clientdesc TEXT NOT NULL UNIQUE,
    clientoff BOOLEAN NOT NULL
);

CREATE TABLE if not exists prj_channel(
    channelid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    channelidtg bigint NOT NULL UNIQUE,
    channelname CHARACTER VARYING(300) NOT NULL UNIQUE,
    channellink TEXT NOT NULL UNIQUE,
    channelclose BOOLEAN NOT NULL
);

CREATE TABLE if not exists prj_session (
    sessionid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    appid CHARACTER VARYING(10) NOT NULL UNIQUE,
    apphash CHARACTER VARYING(100) NOT NULL UNIQUE,
    phonenumber CHARACTER VARYING(16) NOT NULL UNIQUE,
    sessiontxt TEXT NOT NULL
);

CREATE TABLE if not exists ref_client_channel(
    refid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    clientid UUID NOT NULL,
    channelid UUID NOT NULL,
    expirationdate timestamp NOT NULL
);

CREATE TABLE if not exists ref_client_session(
    refid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    clientid UUID NOT NULL,
    sessionid UUID NOT NULL,
    expirationdate timestamp NOT NULL
);
