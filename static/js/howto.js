$(document).ready(function () {

    $('#alert').hide();
    $('#alert_error').hide();

    $('#add').on('submit', function(e) {
        e.preventDefault();

        var id = $('#id').val();
        var serverDestination = $('#serverDestination').val();
        var userDestination = $('#userDestination').val();
        var from = $('#from').val();
        var notes = $('#notes').val();

        $.ajax({
            url: '/addserver',
            type: 'POST',
            data: {
                id: id, 
                serverDestination: serverDestination, 
                userDestination: userDestination, 
                from: from, 
                notes: notes
            },
            success: function(data) {
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                    window.location = "/";
                });
            },
            error: function(data) {
                console.log(data);
                $("#alert_error").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert_error").slideUp(500);
                });
            }
        });

    });

    $('#delete').on('submit', function(e) {
        e.preventDefault();

        var id = $('#id').val();

        $.ajax({
            url: '/deleteserver',
            type: 'POST',
            data: {id: id},
            success: function(data) {
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                    window.location = "/";
                });
            },
            error: function(data) {
                console.log(data);
                $("#alert_error").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert_error").slideUp(500);
                });
            }
        });

    });

    $('#edit').on('submit', function(e) {
        e.preventDefault();

        var id = $('#id').val();
        var serverDestination = $('#serverDestination').val();
        var userDestination = $('#userDestination').val();
        var from = $('#from').val();
        var notes = $('#notes').val();

        $.ajax({
            url: '/editserver',
            type: 'POST',
            data: {
                id: id, 
                serverDestination: serverDestination, 
                userDestination: userDestination, 
                from: from, 
                notes: notes
            },
            success: function(data) {
                $('#interactiontext').val('');
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(700);
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