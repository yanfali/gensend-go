/* global window */
(function(window) {
  'use strict';
  angular.module('generator', []).factory('generatorService', function() {
    var crypto = window.crypto;
    var service = {
      FULL_LOWERCASE: 'abcdefghijklmnopqrstuvwxyz',
      FULL_UPPERCASE: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
      NON_SIMILAR_LOWERCASE: 'abcdefghjkmnpqrtuvwxyz',
      NON_SIMILAR_UPPERCASE: 'ABCDEFGHJKLMNPQRTUVWXYZ',
      STANDARD_NUMBERS: '12346789',
      SIMILAR: 'iIlLoOsS150',
      getRandomUint32: function(length) {
        var array = new Uint32Array(length);
        crypto.getRandomValues(array);
        return array;
      }
    };
    return service;
  });
})(window);
