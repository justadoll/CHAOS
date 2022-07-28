function TakeWebShot(address) {
    Swal.fire({
        title: 'Processing webcam screenshot...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    SendCommand(address, "webshot")
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
            return response.text();
        })
        .then(response => {
            Swal.close();
            window.location.href = 'download/' + response;
        }).catch(err => {
        console.log('Error: ', err);
        Swal.fire({
            icon: 'error',
            title: 'Ops...',
            text: 'Error processing screenshot!',
            footer: err
        });
    });
}