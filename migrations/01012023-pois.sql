-- +migrate Up
CREATE TABLE ccms_schedule (
    channel_name varchar NOT NULL,
    schedule varchar NOT NULL,
    event_type varchar(3) NOT NULL,
    scheduled_date varchar(4) NOT NULL,
    scheduled_time varchar(6) NOT NULL,
    window_start_time varchar(4) NOT NULL,
    window_duration_time varchar(4) NOT NULL,
    break_within_window varchar(3) NOT NULL,
    position_within_break varchar(3) NOT NULL,
    scheduled_length varchar(6) NOT NULL,
    actual_aired_time varchar(6) NOT NULL,
    actual_aired_length varchar(8) NOT NULL,
    actual_aired_position varchar(3) NOT NULL,
    spot_identification varchar(11) NOT NULL,
    status_code varchar(4) NOT NULL,
    user_defined varchar,
    created timestamp DEFAULT current_timestamp,
    updated timestamp DEFAULT current_timestamp,
    PRIMARY KEY(channel_name, scheduled_date, scheduled_time, position_within_break)
);

CREATE TABLE channel_alias (
    channelname varchar NOT NULL,
    aliasname varchar PRIMARY KEY,
    created timestamp DEFAULT current_timestamp,
    updated timestamp DEFAULT current_timestamp
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +migrate StatementEnd

CREATE TRIGGER update_ccms_schedule_timestamp BEFORE UPDATE ON ccms_schedule FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_channel_alias_timestamp BEFORE UPDATE ON channel_alias FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE EXTENSION pg_cron;

GRANT USAGE ON SCHEMA cron TO postgres;

SELECT cron.schedule('1 0 * * *', $$DELETE FROM ccms_schedule WHERE created < NOW() - INTERVAL '24 hours'$$);

-- +migrate Down
SELECT cron.unschedule(1);

DROP EXTENSION pg_cron;

DROP TRIGGER update_channel_alias_timestamp ON channel_alias;

DROP TRIGGER update_ccms_schedule_timestamp ON ccms_schedule;

DROP TABLE channel_alias;

DROP TABLE ccms_schedule;

DROP FUNCTION update_timestamp();

