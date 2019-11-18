## Lokales Bauen und Ausführen

Entweder

-   Binaries im Repository bauen

    ```
    cd cmd/treecli
    go build
    cd ../treeservice
    go build    
    ```
    
    und dann auch
    
-   Binaries lokal ausführen

    ```
    ./cmd/treeservice/treeservice --bind localhost:8090
    ./cmd/treecli/treecli --bind localhost:8091 --remote localhost:8090 [command]  
    ```
    
oder

-   Binaries in den GOPATH installieren

    ```
    go install ./...
    ```
    
    und wenn der GOPATH auch im PATH der Shell hinterlegt ist, einfach
    
-   Binaries direkt ausführen

    ```
    treeservice --bind localhost:8090
    treecli --bind localhost:8091 --remote localhost:8090 [command]  
    ```
        
Außerdem kann durch die Flag `--help` jederzeit eine Hilfe zu verfügbaren Kommandos und Flags ausgegeben werden.    

## Ausführen mit Docker

-   Images bauen

    ```
    make docker
    ```

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    docker run --rm --net actors --name treeservice treeservice \
      --bind="treeservice.actors:8090"
    ```

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    docker run --rm --net actors --name treecli treecli --bind="treecli.actors:8091" \
      --remote="treeservice.actors:8090"
    ```

-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

## Ausführen mit Docker ohne vorher die Docker-Images zu bauen

Nach einem Commit baut der Jenkins, wenn alles durch gelaufen ist, die beiden
Docker-Images. Sie können diese dann mit `docker pull` herunter laden. Schauen Sie für die
genaue Bezeichnung in die Consolenausgabe des Jenkins-Jobs.

Wenn Sie die Imagenamen oben (`treeservice` und `treecli`) durch die Namen aus der
Registry ersetzen, können Sie Ihre Lösung mit den selben Kommandos wie oben beschrieben,
ausprobieren.

## Verfügbare Kommandos

-   Neuen Baum erstellen (MAX_ELEMS ist die Zahl der maximalen Elemente pro Blatt)
    
    ```
    treecli (...) create MAX_ELEMS  
    ```
    
    Gibt die ID und das TOKEN des Baumes zurück, welche für die anderen Operationen gebraucht werden.
    
-   Element einfügen (Schlüssel-Wert Paar KEY-VALUE)

    ```
    treecli (...) --id ID --token TOKEN insert KEY VALUE
    ```
    
-   Element suchen (mit Schlüssel KEY)

    ```
    treecli (...) --id ID --token TOKEN search KEY
    ```
    
    Gibt den Wert des gesuchten Elements zurück
    
-   Element löschen (mit Schlüssel KEY)

    ```
    treecli (...) --id ID --token TOKEN delete KEY
    ```
    
-   Baum traversieren (alle Elemente ausgeben)

    ```
    treecli (...) --id ID --token TOKEN traverse
    ```
    
-   Gesamten Baum entfernen

    ```
    treecli (...) --id ID --token TOKEN remove
    ```
