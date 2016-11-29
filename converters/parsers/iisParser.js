var fieldDefinition = '#Fields:';
function isFieldDefinition(line) {
    return line.slice(0, fieldDefinition.length) == fieldDefinition;
}

function parseLine(line, fields) {

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

    function parseFields(line) {
        line = line.slice(fieldDefinition.length);
        return line.trim().split(' ');
    }

    if (isFieldDefinition(line)) {
        return {fields: parseFields(line)};
    }

    var entry = zip(fields, line.split(' '));

    return JSON.stringify({
        request: {
            method: get(entry, 'cs-method'),
            path: get(entry, 'cs-uri-stem'),
            query: get(entry, 'cs-uri-query'),
            headers: {
                'User-Agent': get(entry, 'cs(User-Agent)'),
                'Referer': get(entry, 'cs(Referer)')
            }
        },
        response: {
            status: parseInt(get(entry, 'sc-status'))
        }
    });
}
