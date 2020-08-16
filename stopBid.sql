CREATE OR REPLACE FUNCTION public.stopbid()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
	declare 
		 temprow record;
         cti bigint :=	extract(epoch from now()) * 1000;
	begin
		for temprow in 
			select b."auctionId" ,b."name" ,max(b.amount) ,b."userId" from auctions a inner join bids b on a.id = b."auctionId" where a."isSold" = false and a."end" < cti group by b."auctionId" ,b."name" ,b."userId"
		loop 
			update auctions set "isActive" = false , "isSold" = true ,"userId" = temprow."userId" where auctions.id = temprow."auctionId" ;
		end loop;
end;
$function$
;
