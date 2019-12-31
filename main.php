<?php

if (! extension_loaded('yaml')) {
    trigger_error('Yaml extension not loaded.', E_USER_ERROR);
}

require_once('ApiDoc.php');

$json = file_get_contents('https://develop.battlenet.com.cn/api/data/navigation/documentation.json');
$nav = json_decode($json);
foreach ($nav->children[0]->children as $child) {
    foreach ($child->page->content->sections as $section) {
        foreach ($section->cardPages as $cardPage) {
            if (property_exists($cardPage, 'contentType') && $cardPage->contentType == 'api-reference') {
                $doc = new ApiDoc();
                $doc->generate($cardPage);
                $doc->publish('battlenet', str_replace('/', '-', substr($cardPage->path, 14)));
            }
        }
    }
}
