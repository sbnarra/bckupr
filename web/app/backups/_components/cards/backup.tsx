import * as spec from '../../../../components/spec'
import {Card, CardBody, Divider} from "@nextui-org/react";

export function BackupCard(props: {
    backup: spec.Backup
}) {
    var backup = props.backup
    return (<Card>
        <CardBody>
            <p>Id: {backup.id}</p>
            <p>Created: {backup.created + ""}</p>
        </CardBody>
        <Divider/>
        <CardBody>
            <p>Status: {backup.status}</p>
            <p>Error: {backup.error}</p>
        </CardBody>
    </Card>)
}