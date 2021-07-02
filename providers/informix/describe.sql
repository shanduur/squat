select upper(c.colname) as colname,
    case
        bitand(255, c.coltype)
        when 0 then 'CHAR'
        when 1 then 'SMALLINT'
        when 2 then 'INTEGER'
        when 3 then 'FLOAT'
        when 4 then 'SMALLFLOAT'
        when 5 then 'DECIMAL'
        when 6 then 'INTEGER' -- serial, autoincrementing int
        when 7 then 'DATE'
        when 10 then case
            c.collength
            when 4365 then 'TIMESTAMP'
            else 'DATETIME'
        end
        when 13 then 'CHAR'
        when 16 then 'CHAR' -- nvarchar
        when 40 then 'CHAR' -- lvarchar
        else cast(concat('UNKNOWN::', to_char(bitand(255, c.coltype))) as varchar)
    end as coltype,
    case
        bitand(255, c.coltype)
        when 5 then trunc(c.collength / 256, 0)
        else c.collength
    end as collength,
    case
        bitand(255, c.coltype)
        when 5 then c.collength - trunc(C.collength / 256, 0) * 256
        else null
    end as precision,
    case
        bitand(c.coltype, 256)
        when 0 then cast('t' as boolean)
        when 256 then cast('f' as boolean)
    end as nullable
from informix.syscolumns c
    inner join informix.systables t on t.tabid = c.tabid
where upper(t.tabname) = upper(?);