import * as spec from '../../../components/spec'
import {Card, CardBody, Divider, Button} from "@nextui-org/react";
import RestoreBackup from './restore_modal';
import DeleteBackup from './delete_modal';

export function BackupCard(props: {
    backup: spec.Backup
}) {
    var backup = props.backup

    return (<Card>
        <CardBody>
            <p>Id: {backup.id}</p>
            <p>Created: {backup.created + ""}</p>
        </CardBody>
        <CardBody>
            <p>Status: {backup.status}</p>
            <p>Error: {backup.error}</p>
        </CardBody>
        <Divider/>
        <CardBody>
            <RestoreBackup backup={backup}/>
            <DeleteBackup backup={backup}/>
        </CardBody>
    </Card>)
}