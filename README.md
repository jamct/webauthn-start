# WebAuthn Demo-Anwendung

Eine Microservice-Anwendung für die Registrierung und Anmeldung mit WebAuthn. Verpackt in Docker-Container und mit einem Docker-Compopse-File versehen.

## Getting started

Zum Hochfahren des Setups reicht der folngende Befehl:

```
docker-compose up
```

## Einsatz in der Praxis

Diese Beispiel-Implementierung ist nicht fertig für den Produktivbetrieb und sollte nur als Inspiration für eigene Dienste genutzt werden!

In den meisten Fällen wollen Sie vermutlich bestehenden Benutzern zusätzlich zum Kennwort (als zweiten Faktor) oder optional (als einzigen Faktor) WebAuthn ermöglichen. Dann können Sie den Container entsprechend umbauen, den Registrierungs-Endpunkt weglassen (die Benutzer liegen ja schon in einer anderen Datenbank) und nur einen neuen Endpoint anbieten, um Token hinzuzufügen.

Einige Probleme sind in dieser Demo noch nicht gelöst: Der Benutzer wird früher oder später seinen Authenticator verlieren. Handelt es sich um ein firmeninternes System kann man damit leben, sofern er einen Administrator kontaktieren kann, der das registrierten Token löscht und dem Nutzer irgendwie ermöglicht, ein neues zu hinterlegen.

Sehen Sie Ihre Nutzer nicht persönlich, braucht es einen Weg, über den sie selbst ihren Token bei Verlust ab- und einen neuen anmelden können. Vergleichsweise einfach zu implementieren ist ein Rücksetzcode, wie ihn auch Google und GitHub anbieten. Generieren Sie ein kurzes Schlüsselpaar, speichern Sie den öffentlichen Schlüssel in der Datenbank und zeigen Sie dem Nutzer den privaten Schlüssel mit der Aufforderung, ihn auszudrucken und an einen sicheren Ort zu legen.

Der Einsatz von WebAuthn in einer produktiven Umgebung ist aber nicht nur ein technisches Problem - auch Ihre Nutzer müssen langsam daran gewöhnt werden, dem Anmelde-Stick mehr zu vertrauen als ihrem Kennwort. Sinnvoll kann es sein, im ersten Schritt optionale Zwei-Faktor-Authentifizierung mit WebAuthn zu aktivieren, diese später verpflichtend zu machen und erst nach etwas Probezeit das Kennwort gänzlich abzuschaffen
