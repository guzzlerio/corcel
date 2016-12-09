module.exports = function parseLine(line, format) {
    function zip(names, values) {
        var result = {};
        for (var i = 0; i < names.length; i++)
            result[names[i]] = values[i];
        return result;
    }

    function get(entry, field) {
        var value = entry[field];
        if (value === '-') {
            return null;
        }
        return value;
    }

    function expandKnownShorthandFormats(format) {
        var request = '"%r"';
        if (format.indexOf(request) > 0) {
            format = format.replace(request, '%m %U%q %H');
        }
        var query = '%U%q';
        if (format.indexOf(query) > 0) {
            format = format.replace(query, '%U %q');
        }
        return format;
    }
    function parseFields(format) {
        var result = [];
        var fields = format.trim().split(' ');
        for (var i = 0; i < fields.length; i++) {
            switch(fields[i]) {
            case '%h':
                result.push('remote-host');
                break;
            case '%l':
                result.push('idenity');
                break;
            case '%u':
                result.push('user');
                break;
            case '%t':
                result.push('time');
                break;
            case '%m':
                result.push('method');
                break;
            case '%U':
                result.push('uri-stem');
                break;
            case '%q':
                result.push('uri-query');
                break;
            case '%H':
                result.push('scheme');
                break;
            case '%>s':
                result.push('status');
                break;
            case '%b':
                result.push('bytes');
                break;
            case '"%{Referer}i"':
                result.push('header-Referer');
                break;
            case '"%{User-agent}i"':
                result.push('header-User-agent');
                break;
            }
        }
        return result;
    }

    // Need to provider the format and the regex as input to this parser.
    // Stanard regex can be provided for CommonFormat and CombinedFormat
    function parseLine(line) {
        var regexLogLine = /^([0-9a-f.:]+) (-) (-) \[([0-9]{2}\/[a-z]{3}\/[0-9]{4}:[0-9]{2}:[0-9]{2}:[0-9]{2}[^\]]*)\] \"(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS) (.+) (HTTP[S]?)\/1.[0-1]\" ([0-9]+) ([0-9]+)/i;
        var parts = regexLogLine.exec(line);
        var a = parts.slice(1, 6);
        a = a.concat(parts[6].split('?'));
        a = a.concat(parts.slice(7));
        return a;
    }

    var fields = parseFields(expandKnownShorthandFormats(format));
    var entry = zip(fields, parseLine(line));

    return JSON.stringify({
        request: {
            method: get(entry, 'method'),
            path: get(entry, 'uri-stem'),
            query: get(entry, 'uri-query'),
            headers: {
                'User-Agent': get(entry, 'header-User-agent'),
                'Referer': get(entry, 'header-Referer')
            }
        },
        response: {
            status: parseInt(get(entry, 'status'))
        }
    });
}
