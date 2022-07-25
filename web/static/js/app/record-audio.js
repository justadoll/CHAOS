async function AudioRecord(address) {
    Swal.fire({
        title: 'Statring microphone record',
        input: 'text',
        reverseButtons: true,
        showCancelButton: true,
        confirmButtonText: 'Start',
        showLoaderOnConfirm: true,
        preConfirm: (url) => {
            Swal.fire({
                title: 'Recoding...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

            sendRecord(address, url)
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(err => {
                            throw new Error(err);
                        });
                    }
                    return response.text();
                })
                .then((result) => {
                    Swal.close();
                    Swal.fire({
                        text: 'Recorded!',
                        icon: 'success'
                    });
                })
                .catch(err => {
                    console.log('Error: ', err);
                    Swal.fire({
                        icon: 'error',
                        title: 'Ops...',
                        text: 'Error recording from microphone!',
                        footer: err
                    });
                })
        },
        allowOutsideClick: () => !Swal.isLoading()
    })
}

async function sendRecord(address, urlToOpen) {
    let formData = new FormData();
    formData.append('address', address);
    formData.append('seconds', urlToOpen);

    const url = '/record-audio';
    const initDetails = {
        method: 'POST',
        body: formData,
        mode: "cors"
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}