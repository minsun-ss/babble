"""create babel tables

Revision ID: c42fec3758d8
Revises:
Create Date: 2024-12-14 18:41:03.915496

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.mysql import BIGINT, LONGBLOB, TIMESTAMP

# revision identifiers, used by Alembic.
revision: str = 'c42fec3758d8'
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        'docs',
        sa.Column('name', sa.String(50), nullable=False, primary_key=True),
        sa.Column('description', sa.String(50), nullable=True),
        sa.Column('hidden', sa.Boolean, default=False),
        sa.Column('last_updated_dt', TIMESTAMP,
            server_default=sa.text('CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP'),
        ),
        sa.Index("ix_last_updated_dt", "last_updated_dt"),
        mysql_engine="InnoDB"
    )

    op.create_table(
        "doc_history",
        sa.Column("id", BIGINT, nullable=False, primary_key=True, autoincrement=True),
        sa.Column('name', sa.String(50), nullable=False),
        sa.Column('version_major', sa.String(10), nullable=False),
        sa.Column("version_minor", sa.String(10), nullable=False),
        sa.Column("version_patch", sa.String(50), nullable=False),
        sa.Column('html', LONGBLOB, nullable=False),
        sa.Column("last_updated_dt", TIMESTAMP, server_default=sa.text('CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP')),
        sa.Index("idx_name", "name"),
        sa.Index("idx_version_major", "version_major"),
        sa.Index("idx_version_minor", "version_minor"),
        sa.Index("idx_version_patch", "version_patch"),
        sa.UniqueConstraint("name", "version_major", "version_minor", "version_patch"),
        mysql_engine="InnoDB"
)

def downgrade() -> None:
    op.drop_table('docs')
    op.drop_table('doc_history')
