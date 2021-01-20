module.exports = {
    'isObj': function (obj) {
        return obj != null && obj.constructor.name == 'Object';
    }
}