$(document).ready(function () {

    $('#alert').hide();
    $('#alert_error').hide();

    if ($('#date').length) {
        $('#date').datepicker({
            format: "yyyy/mm/dd",
            weekStart: 1,
            todayBtn: "linked",
            todayHighlight: true
        });
    }

    $(".today").click();

    var persontypes = $('#persontypes');
    if (persontypes.length) {
        $.ajax({
            url: '/persontypes',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#persontypes').append($('<option name="' + types[i].ID + '">').append(types[i].Type));
                }
            },
            error: function(data) {
                console.log('woops! :(' + data);
            }
        });
    }

    var persons = $('#persons');
    if (persons.length) {
        $.ajax({
            url: '/persons',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#persons').append($('<option name="' + types[i].ID + '">').append(types[i].Name));
                }
            },
            error: function(data) {
                console.log('woops! :(' + data);
            }
        });
    }

    var family = $('#family');
    if (family.length) {
        $.ajax({
            url: '/personspertype/1',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    //$('#family').append($('<li name="1" class="list-group-item">').append(types[i]));
                    $('#family').append(
                        $('<li name="1" class="list-group-item">').append(
                            $('<a>').attr('href','/person/' + types[i].ID).append(
                                $('<span>').attr('class', 'badge badge-light').append(types[i].Name)
                    )));
                }
            },
            error: function(data) {
                console.log('woops! :(');
                console.log(data);
            }
        });
    }

    var friends = $('#friends');
    if (friends.length) {
        $.ajax({
            url: '/personspertype/2',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    //$('#friends').append($('<li name="1" class="list-group-item">').append(types[i]));
                    $('#friends').append(
                        $('<li name="1" class="list-group-item">').append(
                            $('<a>').attr('href','/person/' + types[i].ID).append(
                                $('<span>').attr('class', 'badge badge-light').append(types[i].Name)
                    )));
                }
            },
            error: function(data) {
                console.log('woops! :(');
                console.log(data);
            }
        });
    }

    var coworkers = $('#coworkers');
    if (coworkers.length) {
        $.ajax({
            url: '/personspertype/3',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#coworkers').append(
                        $('<li name="1" class="list-group-item">').append(
                            $('<a>').attr('href','/person/' + types[i].ID).append(
                                $('<span>').attr('class', 'badge badge-light').append(types[i].Person + ' on ' + types[i].Name)
                    )));
                    console.log("Persona: " + types[i]);
                }
            },
            error: function(data) {
                console.log('woops! :(');
                console.log(data);
            }
        });
    }

    $('#addperson').on('submit', function(e) {

        var currentForm = this;
        e.preventDefault();
        var name = $('#person_name').val();
        var personType = $('#persontypes').find(":selected").attr('name');
        var everydays = $('#interacteverydays').find(":selected").val();

        $.ajax({
            url: '/addperson',
            type: 'POST',
            data: {name: name, type: personType, everydays: everydays},
            success: function(data) {
                console.log("Good");
                $('#person_name').val('');
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                });
            },
            error: function(data) {
                console.log("Error!");
                console.log(data);
                $("#alert_error").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert_error").slideUp(500);
                });
            }
        });

    });

    $('#edit').on('submit', function(e) {

        var currentForm = this;
        e.preventDefault();

        var id = $('#id').val();
        var serverDestination = $('#serverDestination').val();
        var userDestination = $('#userDestination').val();
        var from = $('#from').val();
        var notes = $('#notes').val();

        $.ajax({
            url: '/editserver',
            type: 'POST',
            data: {id: id, serverDestination: serverDestination, userDestination: userDestination, from: from, notes: notes},
            success: function(data) {
                $('#interactiontext').val('');
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                    window.location = "/";
                });
            },
            error: function(data) {
                $("#alert_error").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert_error").slideUp(500);
                });
            }
        });

    });

});