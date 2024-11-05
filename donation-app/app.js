document.getElementById("donationForm").onsubmit = function(event) {
    event.preventDefault();
    let amount = document.getElementById("amount").value;

    if (!amount || amount <= 0) {
        alert("Jumlah donasi harus lebih besar dari 0.");
        return;
    }

    fetch("/donate", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: "amount=" + encodeURIComponent(amount),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Jaringan bermasalah atau server mengembalikan kesalahan');
        }
        return response.json();
    })
    .then(data => {
        alert(data.message);
        document.getElementById("amount").value = "";
    })
    .catch(error => {
        alert("Terjadi kesalahan: " + error.message);
    });
};
